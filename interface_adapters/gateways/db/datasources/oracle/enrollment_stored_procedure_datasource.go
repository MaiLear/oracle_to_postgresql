package oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"github.com/godror/godror"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/types"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/logger"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
)

type EnrollmentStoredProcedureDataSource struct {
	connection *sql.DB
}

func NewEnrollmentStoredProcedureDataSource(connection *sql.DB) EnrollmentStoredProcedureDataSource {
	return EnrollmentStoredProcedureDataSource{
		connection: connection,
	}
}

// EnrollmentStoredProcedureParams contiene todos los parámetros necesarios para el procedimiento
type EnrollmentStoredProcedureParams struct {
	DocumentType         string
	DocumentNumber       string
	Password             string
	CourseId             int
	PopulationTypeId     *int
	AgreementNumber      *string
	SurveyResponses      []SurveyResponse
	IcfesCodes           []IcfesCode
	SecondOptionCourseId *int
	SpecialOfferTypeId   *int
}

// SurveyResponse representa una respuesta de encuesta para el procedimiento
type SurveyResponse struct {
	RinId int
	CnoId int
}

// IcfesCode representa un código ICFES para el procedimiento
type IcfesCode struct {
	SnpCode        string
	DocumentType   string
	DocumentNumber string
}

// Enroll ejecuta el procedimiento almacenado USP_INSERTAR_INSCRIPCION2
func (esp EnrollmentStoredProcedureDataSource) Enroll(ctx context.Context, enrollmentDomain entities.Enrollment) (string,error) {
	if enrollmentDomain.RequestData == nil{
		return "",cockroachdbErrors.New("infra: el campo RequetData es nil por lo que no se puede obtener la informacion para insertar")
	}
	params := *enrollmentDomain.RequestData
	logger.InfoLogger.Printf("🎓 Ejecutando procedimiento USP_INSERTAR_INSCRIPCION2 para curso: %d, documento: %s - %s",
		params.CourseID, params.DocumentType, params.DocumentNumber)
	logger.InfoLogger.Printf("📊 Parámetros: SurveyResponses=%d, IcfesCodes=%d", len(params.SurveyResponses), len(params.IcfesCodes))

	// Log detallado de todos los parámetros
	logger.InfoLogger.Printf("🔍 PARÁMETROS DETALLADOS DEL PROCEDIMIENTO:")
	logger.InfoLogger.Printf("   - DocumentType: %s", params.DocumentType)
	logger.InfoLogger.Printf("   - DocumentNumber: %s", params.DocumentNumber)
	logger.InfoLogger.Printf("   - CourseId: %d", params.CourseID)
	logger.InfoLogger.Printf("   - PopulationTypeId: %v", params.PopulationTypeID)
	logger.InfoLogger.Printf("   - AgreementNumber: %v", params.AgreementNumber)
	logger.InfoLogger.Printf("   - SecondOptionCourseId: %v", params.SecondOptionCourseID)
	logger.InfoLogger.Printf("   - SpecialOfferTypeId: %v", params.SpecialOfferTypeID)
	logger.InfoLogger.Printf("   - CONTRASEÑA: %v", params.Password)

	// Log de respuestas de encuesta
	if len(params.SurveyResponses) > 0 {
		logger.InfoLogger.Printf("   - SurveyResponses:")
		for i, response := range params.SurveyResponses {
			logger.InfoLogger.Printf("     [%d] RinId: %d, CnoId: %d", i+1, response.RinId, response.CnoId)
		}
	} else {
		logger.InfoLogger.Printf("   - SurveyResponses: [] (vacío)")
	}

	// Log de códigos ICFES
	if len(params.IcfesCodes) > 0 {
		logger.InfoLogger.Printf("   - IcfesCodes:")
		for i, code := range params.IcfesCodes {
			logger.InfoLogger.Printf("     [%d] SnpCode: %s, DocumentType: %s, DocumentNumber: %s",
				i+1, code.SnpCode, code.DocumentType, code.DocumentNumber)
		}
	} else {
		logger.InfoLogger.Printf("   - IcfesCodes: [] (vacío)")
	}

	// Obtener una conexión de godror para manejar tipos personalizados
	conn, err := esp.connection.Conn(ctx)
	if err != nil {
		return "",cockroachdbErrors.Wrapf(err, "infra: error obteniendo conexión para tipos personalizados: %s",err.Error())
	}
	defer conn.Close()

	// Crear los arrays personalizados de Oracle usando godror
	surveyResponsesArray, err := esp.createSurveyResponsesArray(ctx, conn, params.SurveyResponses)
	if err != nil {
		return "",cockroachdbErrors.Wrapf(err, "infra: error creando array de respuestas de encuesta: %s",err.Error())
	}
	logger.InfoLogger.Printf("✅ Array de respuestas creado: %v", surveyResponsesArray)

	icfesCodesArray, err := esp.createIcfesCodesArray(ctx, conn, params.IcfesCodes)
	if err != nil {
		return "",cockroachdbErrors.Wrapf(err, "infra: error creando array de códigos ICFES: %s",err.Error())
	}
	logger.InfoLogger.Printf("✅ Array de códigos ICFES creado: %v", icfesCodesArray)

	// Construir el SQL para llamar al procedimiento
	callSQL := `BEGIN
		INSCRIPCION.USP_INSERTAR_INSCRIPCION2(
			:tipo_id, :num_id, :clave, :fic_id, :id_tipo_pob, :cov_numero, :respuestas, :snp, :fic_id_segunda_opc, :toe_id, :resultado
		);
	END;`

	// Variables para valores opcionales
	var populationTypeId sql.NullInt64
	if params.PopulationTypeID != nil {
		populationTypeId.Int64 = int64(*params.PopulationTypeID)
		populationTypeId.Valid = true
	}

	var agreementNumber sql.NullString
	if params.AgreementNumber != nil {
		agreementNumber.String = *params.AgreementNumber
		agreementNumber.Valid = true
	}

	var secondOptionCourseId sql.NullInt64
	if params.SecondOptionCourseID != nil {
		secondOptionCourseId.Int64 = *params.SecondOptionCourseID
		secondOptionCourseId.Valid = true
	}

	var specialOfferTypeId sql.NullInt64
	if params.SpecialOfferTypeID != nil {
		specialOfferTypeId.Int64 = int64(*params.SpecialOfferTypeID)
		specialOfferTypeId.Valid = true
	}

	// Variable para capturar el resultado
	var resultado string

	// Ejecutar el procedimiento usando la conexión de godror
	_, err = conn.ExecContext(ctx, callSQL,
		sql.Named("tipo_id", params.DocumentType),
		sql.Named("num_id", params.DocumentNumber),
		//Cris debe agregar la contraseña al obtejeto de datos
		sql.Named("clave", params.Password),
		sql.Named("fic_id", params.CourseID),
		sql.Named("id_tipo_pob", populationTypeId),
		sql.Named("cov_numero", agreementNumber),
		sql.Named("respuestas", surveyResponsesArray),
		sql.Named("snp", icfesCodesArray),
		sql.Named("fic_id_segunda_opc", secondOptionCourseId),
		sql.Named("toe_id", specialOfferTypeId),
		sql.Named("resultado", sql.Out{Dest: &resultado}),
	)

	if err != nil {
		logger.ErrorLogger.Printf("❌ ERROR EJECUTANDO PROCEDIMIENTO USP_INSERTAR_INSCRIPCION2: %v", err)
		// Manejar errores específicos de Oracle
		errStr := err.Error()
		return errStr,cockroachdbErrors.Wrapf(err,"infra: error ejecutando procedimiento de inscripción: %s",errStr)

		// // Error de aplicación (mensaje para el usuario)
		// if strings.Contains(errStr, "ORA-20001") {
		// 	// Extraer el mensaje de error de la aplicación
		// 	errorMessage := extractApplicationError(errStr)
		// 	return errorMessage,cockroachdbErrors.Wrapf(err,"infra: %s",errorMessage)
		// 	// return "", &internalErrors.AppError{
		// 	// 	Message: errorMessage,
		// 	// 	Err:     err,
		// 	// 	Details: map[string]any{
		// 	// 		"course_id":               params.CourseId,
		// 	// 		"document_type":           params.DocumentType,
		// 	// 		"document_number":         params.DocumentNumber,
		// 	// 		"population_type_id":      params.PopulationTypeId,
		// 	// 		"second_option_course_id": params.SecondOptionCourseId,
		// 	// 		"special_offer_type_id":   params.SpecialOfferTypeId,
		// 	// 		"source":                  "ORACLE_PROCEDURE",
		// 	// 		"oracle_error":            errStr,
		// 	// 	},
		// 	// }
		// }

		// // Error de restricción única (usuario ya inscrito)
		// if strings.Contains(errStr, "ORA-00001") && strings.Contains(errStr, "ING_NIS_FIC_UK") {
		// 	errMessage := "El usuario ya se encuentra inscrito en este programa de formación"
		// 	return errMessage,cockroachdbErrors.Wrapf(internalErrors.ErrViolationUnique,"%s",errMessage)
		// 	// return "", &internalErrors.AppError{
		// 	// 	Message: "El usuario ya se encuentra inscrito en este programa de formación",
		// 	// 	Err:     err,
		// 	// 	Details: map[string]any{
		// 	// 		"course_id":            params.CourseId,
		// 	// 		"document_type":        params.DocumentType,
		// 	// 		"document_number":      params.DocumentNumber,
		// 	// 		"source":               "ORACLE_PROCEDURE",
		// 	// 		"oracle_error":         errStr,
		// 	// 		"constraint_violation": "ING_NIS_FIC_UK",
		// 	// 	},
		// 	// }
		// }

		// // Error de datos no encontrados (usuario no existe o datos faltantes)
		// if strings.Contains(errStr, "ORA-01403") {
		// 	// Determinar el contexto específico del error basado en el procedimiento
		// 	var errorMessage string
		// 	if strings.Contains(errStr, "CONFIRMAR_INS_VIRTUAL") {
		// 		errorMessage = "No se pudo completar la confirmación de la inscripción. Por favor, intente inscribirse nuevamente o contacte al soporte técnico si el problema persiste."
		// 		return errorMessage,cockroachdbErrors.Wrapf(internalErrors.ErrBusiness,"infra: %s",errorMessage)
		// 	} else {
		// 		errorMessage = "no se encontraron los datos necesarios para completar la inscripción. Por favor, verifique que su información personal esté completa y actualizada en el sistema SOFIA Plus."
		// 		return errorMessage,cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,"infra: %s",errorMessage)
		// 	}

		// 	// return "", &internalErrors.AppError{
		// 	// 	Message: errorMessage,
		// 	// 	Err:     err,
		// 	// 	Details: map[string]any{
		// 	// 		"course_id":       params.CourseId,
		// 	// 		"document_type":   params.DocumentType,
		// 	// 		"document_number": params.DocumentNumber,
		// 	// 		"source":          "ORACLE_PROCEDURE",
		// 	// 		"oracle_error":    errStr,
		// 	// 		"error_type":      "NO_DATA_FOUND",
		// 	// 		"procedure":       "CONFIRMAR_INS_VIRTUAL",
		// 	// 	},
		// 	// }
		// }

		// Error de base de datos (log interno)
		// logger.ErrorLogger.Printf("❌ Error ejecutando procedimiento USP_INSERTAR_INSCRIPCION2: %v", err)
		// return "",cockroachdbErrors.Wrap(err, "infra: error ejecutando procedimiento de inscripción")
	}

	// Verificar si el resultado contiene un mensaje de error o validación
	// Palabras clave más específicas para detectar errores reales
	errorKeywords := []string{
		"error", "no se pudo", "no es válido", "no es valido", "no existe",
		"ya existe", "no encontrado", "no se encuentra disponible",
		"seleccione otra", "inhabilitado", "inscrito",
		"no se puede realizar", "no se encuentra", "no es correcta",
		"no se encuentra disponible una ficha", "no existen más niveles",
		"no se puede realizar la inscripción", "ya cuenta con un certificado",
		"por certificar", "estado por certificar", "no puede realizar",
		"no es correcta", "inconsistencia en datos", "no se ha definido",
		"no se encuentra registrado", "no se encuentra habilitado",
		"no puede validar", "debe transcurrir", "debe adjuntar",
		"presenta una novedad", "ha sido reportada", "no ha sido tramitada",
		"tiene registrada", "existe una sancion", "no existen fechas activas",
		"no puede verificar", "no puede continuar", "debe acercarse",
		"debe contactarse", "debe actualizar", "debe corregir",
	}

	isError := false
	for _, keyword := range errorKeywords {
		if strings.Contains(strings.ToLower(resultado), strings.ToLower(keyword)) {
			isError = true
			break
		}
	}

	// Verificar si es un mensaje de éxito (contiene palabras específicas de éxito)
	successKeywords := []string{
		"se ha inscrito satisfactoriamente", "inscrito satisfactoriamente",
		"cuenta de correo", "sofia plus", "registro persona",
	}

	isSuccess := false
	for _, keyword := range successKeywords {
		if strings.Contains(strings.ToLower(resultado), strings.ToLower(keyword)) {
			isSuccess = true
			break
		}
	}

	// Si es un mensaje de éxito, retornarlo directamente sin verificar errores
	if isSuccess {
		logger.InfoLogger.Printf("✅ PROCEDIMIENTO EXITOSO: %s", resultado)
		return "",nil
	}

	// Solo verificar errores si NO es un mensaje de éxito
	if isError {
		// Es un mensaje de error real, retornarlo como AppError
		errMessage := cleanErrorMessage(resultado)
		return errMessage,cockroachdbErrors.Wrapf(internalErrors.ErrBusiness,"infra: %s",resultado,)
		// return "", &internalErrors.AppError{
		// 	Message: cleanErrorMessage(resultado),
		// 	Err:     errors.New("error en procedimiento de inscripción"),
		// 	Details: map[string]any{
		// 		"course_id":               params.CourseId,
		// 		"document_type":           params.DocumentType,
		// 		"document_number":         params.DocumentNumber,
		// 		"population_type_id":      params.PopulationTypeId,
		// 		"second_option_course_id": params.SecondOptionCourseId,
		// 		"special_offer_type_id":   params.SpecialOfferTypeId,
		// 		"source":                  "ORACLE_PROCEDURE",
		// 	},
		// }
	}

	// Si no es error ni éxito, asumir que es éxito (comportamiento por defecto)
	logger.InfoLogger.Printf("✅ Procedimiento USP_INSERTAR_INSCRIPCION2 ejecutado exitosamente.")
	logger.InfoLogger.Printf("📋 RESULTADO DEL PROCEDIMIENTO: %s", resultado)
	logger.InfoLogger.Printf("🔚 FIN EJECUCIÓN PROCEDIMIENTO USP_INSERTAR_INSCRIPCION2")
	return  "",nil
}

// EnrollComplementary ejecuta el procedimiento almacenado PR_INS_INSCRIPCION_VIRTUAL para programas complementarios
func (esp EnrollmentStoredProcedureDataSource) EnrollComplementary(ctx context.Context, enrollmentDomain entities.Enrollment) (string,error) {
	if enrollmentDomain.RequestData == nil{
		return "",cockroachdbErrors.New("infra: el campo RequetData es nil por lo que no se puede obtener la informacion para insertar")
	}
	params := *enrollmentDomain.RequestData
	logger.InfoLogger.Printf("🎓 Ejecutando procedimiento PR_INS_INSCRIPCION_VIRTUAL para NIS: %d, FIC_ID: %d", enrollmentDomain.Nis, params.CourseID)
	logger.InfoLogger.Printf("📊 Parámetros: AgreementNumber=%v", params.AgreementNumber)

	// Log detallado de todos los parámetros (como en la plataforma Java)
	logger.InfoLogger.Printf("🔍 PARÁMETROS DETALLADOS DEL PROCEDIMIENTO COMPLEMENTARIO:")
	logger.InfoLogger.Printf("   - p_nis: %d", params.Nis)
	logger.InfoLogger.Printf("   - p_fic_id: %d", params.CourseID)
	logger.InfoLogger.Printf("   - p_cov_numero: %v", params.AgreementNumber)
	logger.InfoLogger.Printf("   - p_nis_fun_registro: '' (valor por defecto)")
	logger.InfoLogger.Printf("   - p_es_familia: '0' (valor por defecto)")
	logger.InfoLogger.Printf("   - p_centro: '' (valor por defecto)")
	logger.InfoLogger.Printf("   - p_resultado: OUT")
	logger.InfoLogger.Printf("   - p_ficha_familia: OUT")

	// Obtener una conexión de godror para manejar tipos personalizados
	conn, err := esp.connection.Conn(ctx)
	if err != nil {
		return "",cockroachdbErrors.Wrapf(err, "infra: error obteniendo conexión para tipos personalizados: %s",err.Error())
	}
	defer conn.Close()

	// Crear los arrays personalizados de Oracle usando godror
	surveyResponsesArray, err := esp.createSurveyResponsesArray(ctx, conn, params.SurveyResponses)
	if err != nil {
		return "",cockroachdbErrors.Wrapf(err, "infra: error creando array de respuestas de encuesta: %s",err.Error())
	}
	logger.InfoLogger.Printf("✅ Array de respuestas creado: %v", surveyResponsesArray)

	icfesCodesArray, err := esp.createIcfesCodesArray(ctx, conn, params.IcfesCodes)
	if err != nil {
		return "",cockroachdbErrors.Wrapf(err, "infra: error creando array de códigos ICFES: %s",err.Error())
	}
	logger.InfoLogger.Printf("✅ Array de códigos ICFES creado: %v", icfesCodesArray)

	// Construir el SQL para llamar al procedimiento complementario
	callSQL := `BEGIN
		INSCRIPCION.PR_INS_INSCRIPCION_VIRTUAL(
			:p_nis, :p_fic_id, :p_cov_numero, :p_nis_fun_registro, :p_es_familia, :p_centro, :p_resultado, :p_ficha_familia
		);
	END;`

	// Variables para valores opcionales
	var covNumero sql.NullString
	if params.AgreementNumber != nil {
		covNumero.String = *params.AgreementNumber
		covNumero.Valid = true
	}

	// Variables para capturar los resultados
	var resultado string
	var fichaFamilia sql.NullInt64

	// NIS ya viene como parámetro

	// Log específico como en la plataforma Java
	logger.InfoLogger.Printf("Saliendo de arregloDeParametros Inscripción Virtual. valores de Parámetros : (IN:'%d', IN:'%d', IN:'%v', IN:'',IN:'0',IN:'',OUT:'',OUT:'')",
		params.Nis, params.CourseID, params.AgreementNumber)

	// Log adicional para diagnosticar el problema
	logger.InfoLogger.Printf("🔍 DIAGNÓSTICO - Antes de ejecutar PR_INS_INSCRIPCION_VIRTUAL:")
	logger.InfoLogger.Printf("   - NIS: %d", params.Nis)
	logger.InfoLogger.Printf("   - FIC_ID: %d", params.CourseID)
	logger.InfoLogger.Printf("   - AgreementNumber: %v", params.AgreementNumber)
	logger.InfoLogger.Printf("   - DocumentType: %s", params.DocumentType)
	logger.InfoLogger.Printf("   - DocumentNumber: %s", params.DocumentNumber)

	// Verificar que el usuario existe en la base de datos antes de ejecutar el procedimiento
	logger.InfoLogger.Printf("🔍 Verificando existencia del usuario en la base de datos...")
	userExists, err := esp.verifyUserExists(ctx, conn, params.Nis, params.DocumentType, params.DocumentNumber)
	if err != nil {
		logger.ErrorLogger.Printf("❌ Error verificando usuario: %v", err)
	} else {
		logger.InfoLogger.Printf("✅ Usuario existe en la base de datos: %t", userExists)
	}

	// Ejecutar el procedimiento usando la conexión de godror
	_, err = conn.ExecContext(ctx, callSQL,
		sql.Named("p_nis", params.Nis),
		sql.Named("p_fic_id", params.CourseID),
		sql.Named("p_cov_numero", covNumero),
		sql.Named("p_nis_fun_registro", ""), // Valor por defecto como en Java
		sql.Named("p_es_familia", "0"),      // Valor por defecto como en Java
		sql.Named("p_centro", ""),           // Valor por defecto como en Java
		sql.Named("p_resultado", sql.Out{Dest: &resultado}),
		sql.Named("p_ficha_familia", sql.Out{Dest: &fichaFamilia}),
	)

	if err != nil {
		logger.ErrorLogger.Printf("❌ ERROR EJECUTANDO PROCEDIMIENTO PR_INS_INSCRIPCION_VIRTUAL: %v", err)
		// Manejar errores específicos de Oracle
		errStr := err.Error()
		return errStr,cockroachdbErrors.Wrapf(err, "infra: error ejecutando procedimiento de inscripción complementaria: %s",errStr)

		// // Error de aplicación (mensaje para el usuario)
		// if strings.Contains(errStr, "ORA-20001") {
		// 	// Extraer el mensaje de error de la aplicación
		// 	errorMessage := extractApplicationError(errStr)
		// 	return errorMessage,cockroachdbErrors.Wrapf(err,"infra: %s",errorMessage)
		// 	// return "", &internalErrors.AppError{
		// 	// 	Message: errorMessage,
		// 	// 	Err:     err,
		// 	// 	Details: map[string]any{
		// 	// 		"nis":                     nis,
		// 	// 		"course_id":               params.CourseId,
		// 	// 		"document_type":           params.DocumentType,
		// 	// 		"document_number":         params.DocumentNumber,
		// 	// 		"population_type_id":      params.PopulationTypeId,
		// 	// 		"second_option_course_id": params.SecondOptionCourseId,
		// 	// 		"special_offer_type_id":   params.SpecialOfferTypeId,
		// 	// 		"source":                  "ORACLE_PROCEDURE_COMPLEMENTARY",
		// 	// 		"oracle_error":            errStr,
		// 	// 	},
		// 	// }
		// }

		// // Error de restricción única (usuario ya inscrito)
		// if strings.Contains(errStr, "ORA-00001") && strings.Contains(errStr, "ING_NIS_FIC_UK") {
		// 	errMessage := "El usuario ya se encuentra inscrito en este programa de formación"
		// 	return errMessage,cockroachdbErrors.Wrapf(internalErrors.ErrViolationUnique,"infra: %s",errMessage)
		// 	// return "", &internalErrors.AppError{
		// 	// 	Message: "El usuario ya se encuentra inscrito en este programa de formación",
		// 	// 	Err:     err,
		// 	// 	Details: map[string]any{
		// 	// 		"nis":                  nis,
		// 	// 		"course_id":            params.CourseId,
		// 	// 		"document_type":        params.DocumentType,
		// 	// 		"document_number":      params.DocumentNumber,
		// 	// 		"source":               "ORACLE_PROCEDURE_COMPLEMENTARY",
		// 	// 		"oracle_error":         errStr,
		// 	// 		"constraint_violation": "ING_NIS_FIC_UK",
		// 	// 	},
		// 	// }
		// }

		// // Error de datos no encontrados (usuario no existe o datos faltantes)
		// if strings.Contains(errStr, "ORA-01403") {
		// 	// Determinar el contexto específico del error basado en el procedimiento
		// 	var errorMessage string
		// 	//var errorDetails string

		// 	if strings.Contains(errStr, "CONFIRMAR_INS_VIRTUAL") {
		// 		errorMessage = "No se pudo completar la confirmación de la inscripción. Por favor, intente inscribirse nuevamente o contacte al soporte técnico si el problema persiste."
		// 		return errorMessage,cockroachdbErrors.Wrapf(internalErrors.ErrBusiness,"infra: %s",errorMessage)
		// 		//errorDetails = "enrollment_confirmation_failed"
		// 	} else {
		// 		errorMessage = "No se encontraron los datos necesarios para completar la inscripción. Por favor, verifique que su información personal esté completa y actualizada en el sistema SOFIA Plus."
		// 		return errorMessage,cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,"infra: %s",errorMessage)
		// 		//errorDetails = "user_data_not_found"
		// 	}

		// 	// return "", &internalErrors.AppError{
		// 	// 	Message: errorMessage,
		// 	// 	Err:     err,
		// 	// 	Details: map[string]any{
		// 	// 		"nis":             enrollmentDomain.Nis,
		// 	// 		"course_id":       params.CourseID,
		// 	// 		"document_type":   params.DocumentType,
		// 	// 		"document_number": params.DocumentNumber,
		// 	// 		"source":          "ORACLE_PROCEDURE_COMPLEMENTARY",
		// 	// 		"oracle_error":    errStr,
		// 	// 		"error_type":      "NO_DATA_FOUND",
		// 	// 		"procedure":       "CONFIRMAR_INS_VIRTUAL",
		// 	// 		"error_details":   errorDetails,
		// 	// 	},
		// 	// }
		// }

		// // Error de base de datos (log interno)
		// logger.ErrorLogger.Printf("❌ Error ejecutando procedimiento PR_INS_INSCRIPCION_VIRTUAL: %v", err)
		// return "",cockroachdbErrors.Wrap(err, "infra: error ejecutando procedimiento de inscripción complementaria")
	}

	// Log de resultados para diagnóstico
	logger.InfoLogger.Printf("🔍 DIAGNÓSTICO - Después de ejecutar PR_INS_INSCRIPCION_VIRTUAL:")
	logger.InfoLogger.Printf("   - Resultado: %s", resultado)
	logger.InfoLogger.Printf("   - FichaFamilia: %v", fichaFamilia)
	logger.InfoLogger.Printf("   - Error: %v", err)

	// Verificar si hay discrepancia entre FIC_ID y FichaFamilia
	if fichaFamilia.Valid && fichaFamilia.Int64 != int64(params.CourseID) {
		logger.ErrorLogger.Printf("⚠️  ADVERTENCIA: Discrepancia entre FIC_ID y FichaFamilia:")
		logger.ErrorLogger.Printf("   - FIC_ID enviado: %d", params.CourseID)
		logger.ErrorLogger.Printf("   - FichaFamilia retornada: %d", fichaFamilia.Int64)
		logger.ErrorLogger.Printf("   - Esta discrepancia puede causar el error ORA-01403 en CONFIRMAR_INS_VIRTUAL")
	}

	// Verificar si la inscripción se creó correctamente
	logger.InfoLogger.Printf("🔍 Verificando si la inscripción se creó correctamente...")
	enrollmentExists, err := esp.verifyEnrollmentExists(ctx, conn, enrollmentDomain.Nis, params.CourseID)
	if err != nil {
		logger.ErrorLogger.Printf("❌ Error verificando inscripción: %v", err)
	} else {
		logger.InfoLogger.Printf("✅ Inscripción existe en la base de datos: %t", enrollmentExists)
	}

	// Verificar si el resultado contiene un mensaje de error o validación
	// Palabras clave más específicas para detectar errores reales
	errorKeywords := []string{
		"errorval001", "no se pudo", "no es válido", "no es valido", "no existe",
		"ya existe", "no encontrado", "no se encuentra disponible",
		"seleccione otra", "inhabilitado", "inscrito",
		"no se puede realizar", "no se encuentra", "no es correcta",
		"no se encuentra disponible una ficha", "no existen más niveles",
		"no se puede realizar la inscripción", "ya cuenta con un certificado",
	}

	isError := false
	for _, keyword := range errorKeywords {
		if strings.Contains(strings.ToLower(resultado), strings.ToLower(keyword)) {
			isError = true
			break
		}
	}

	// Verificar si es un mensaje de éxito (contiene palabras específicas de éxito)
	successKeywords := []string{
		"se ha inscrito satisfactoriamente", "inscrito satisfactoriamente",
		"cuenta de correo", "sofia plus", "registro persona",
		"se ha inscrito satisfactoriamente y debe presentar",
		"se ha inscrito satisfactoriamente y que debe presentar",
		"se ha inscrito satisfactoriamente y que el proceso",
		"se ha inscrito satisfactoriamente. el sistema le enviará",
		"se ha inscrito satisfactoriamente. el sistema le enviara",
		"quedo inscrito a las",
	}

	isSuccess := false
	for _, keyword := range successKeywords {
		if strings.Contains(strings.ToLower(resultado), strings.ToLower(keyword)) {
			isSuccess = true
			break
		}
	}

	// Si es un mensaje de éxito, retornarlo directamente sin verificar errores
	if isSuccess {
		logger.InfoLogger.Printf("✅ PROCEDIMIENTO COMPLEMENTARIO EXITOSO: %s", resultado)
		return "",nil
	}

	// Solo verificar errores si NO es un mensaje de éxito
	if isError {
		errMessage := cleanErrorMessage(resultado)
		return errMessage,cockroachdbErrors.Wrapf(internalErrors.ErrBusiness,"infra: %s",resultado)
		// Es un mensaje de error real, retornarlo como AppError
		// return "", &internalErrors.AppError{
		// 	Message: cleanErrorMessage(resultado),
		// 	Err:     errors.New("error en procedimiento de inscripción complementaria"),
		// 	Details: map[string]any{
		// 		"nis":                     nis,
		// 		"course_id":               params.CourseId,
		// 		"document_type":           params.DocumentType,
		// 		"document_number":         params.DocumentNumber,
		// 		"population_type_id":      params.PopulationTypeId,
		// 		"second_option_course_id": params.SecondOptionCourseId,
		// 		"special_offer_type_id":   params.SpecialOfferTypeId,
		// 		"source":                  "ORACLE_PROCEDURE_COMPLEMENTARY",
		// 	},
		// }
	}

	// Si no es error ni éxito, asumir que es éxito (comportamiento por defecto)
	logger.InfoLogger.Printf("✅ Procedimiento PR_INS_INSCRIPCION_VIRTUAL ejecutado exitosamente.")
	logger.InfoLogger.Printf("📋 RESULTADO DEL PROCEDIMIENTO COMPLEMENTARIO: %s", resultado)
	logger.InfoLogger.Printf("🔚 FIN EJECUCIÓN PROCEDIMIENTO PR_INS_INSCRIPCION_VIRTUAL")
	return "",nil
}

// createSurveyResponsesArray crea un array personalizado de Oracle para las respuestas de encuesta
// usando godror para manejar el tipo INSCRIPCION.ARR_RESPUESTAS_ENC
func (esp EnrollmentStoredProcedureDataSource) createSurveyResponsesArray(ctx context.Context, conn *sql.Conn, responses []types.SurveyResponse) (interface{}, error) {
	// Obtener el tipo de array Oracle para INSCRIPCION.ARR_RESPUESTAS_ENC
	arrType, err := godror.GetObjectType(ctx, conn, "INSCRIPCION.ARR_RESPUESTAS_ENC")
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "error obteniendo tipo de array ARR_RESPUESTAS_ENC")
	}

	// Crear el array de objetos (siempre crear el array, aunque esté vacío)
	arrayObj, err := arrType.NewObject()
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "error creando objeto de array")
	}

	if len(responses) == 0 {
		// Retornar array vacío en lugar de nil
		return arrayObj, nil
	}

	// Obtener el tipo de objeto Oracle para INSCRIPCION.OBJ_RESPUESTAS_ENC
	objType, err := godror.GetObjectType(ctx, conn, "INSCRIPCION.OBJ_RESPUESTAS_ENC")
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "error obteniendo tipo de objeto OBJ_RESPUESTAS_ENC")
	}

	// Poblar el array con los objetos individuales
	for _, response := range responses {
		obj, err := objType.NewObject()
		if err != nil {
			return nil, cockroachdbErrors.Wrap(err, "error creando objeto de respuesta")
		}

		// Establecer los valores del objeto
		if err := obj.Set("RIN_ID", int64(response.RinId)); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "error estableciendo RIN_ID")
		}
		if err := obj.Set("CNO_ID", fmt.Sprintf("%d", response.CnoId)); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "error estableciendo CNO_ID")
		}

		// Agregar el objeto al array
		if err := arrayObj.Collection().Append(obj); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "error agregando objeto al array")
		}
	}

	return arrayObj, nil
}

// createIcfesCodesArray crea un array personalizado de Oracle para los códigos ICFES
// usando godror para manejar el tipo INSCRIPCION.ARR_SNPXINSCRIPCION
func (esp EnrollmentStoredProcedureDataSource) createIcfesCodesArray(ctx context.Context, conn *sql.Conn, codes []types.Icfes) (interface{}, error) {
	// Obtener el tipo de array Oracle para INSCRIPCION.ARR_SNPXINSCRIPCION
	arrType, err := godror.GetObjectType(ctx, conn, "INSCRIPCION.ARR_SNPXINSCRIPCION")
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "error obteniendo tipo de array ARR_SNPXINSCRIPCION")
	}

	// Crear el array de objetos (siempre crear el array, aunque esté vacío)
	arrayObj, err := arrType.NewObject()
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "error creando objeto de array")
	}

	if len(codes) == 0 {
		// Retornar array vacío en lugar de nil
		return arrayObj, nil
	}

	// Obtener el tipo de objeto Oracle para INSCRIPCION.OBJ_SNPXINSCRIPCION
	objType, err := godror.GetObjectType(ctx, conn, "INSCRIPCION.OBJ_SNPXINSCRIPCION")
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "error obteniendo tipo de objeto OBJ_SNPXINSCRIPCION")
	}

	// Poblar el array con los objetos individuales
	for _, code := range codes {
		obj, err := objType.NewObject()
		if err != nil {
			return nil, cockroachdbErrors.Wrap(err, "error creando objeto de código ICFES")
		}

		// Establecer los valores del objeto
		if err := obj.Set("SNP_CODIGO", code.SnpCode); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "error estableciendo SNP_CODIGO")
		}
		if err := obj.Set("TIPO_DOCUMENTO", code.DocumentType); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "error estableciendo TIPO_DOCUMENTO")
		}
		if err := obj.Set("NUM_DOC_IDENTIDAD", code.DocumentNumber); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "error estableciendo NUM_DOC_IDENTIDAD")
		}

		// Agregar el objeto al array
		if err := arrayObj.Collection().Append(obj); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "error agregando objeto al array")
		}
	}

	return arrayObj, nil
}

// extractApplicationError extrae el mensaje de error de aplicación del error de Oracle
func extractApplicationError(errStr string) string {
	// Buscar el mensaje después de "ORA-20001"
	if strings.Contains(errStr, "ORA-20001") {
		parts := strings.Split(errStr, "ORA-20001")
		if len(parts) > 1 {
			// Limpiar el mensaje de error
			message := strings.TrimSpace(parts[1])
			// Remover información adicional de Oracle
			if strings.Contains(message, "-ERROR-") {
				message = strings.Split(message, "-ERROR-")[0]
			}
			return cleanErrorMessage(strings.TrimSpace(message))
		}
	}
	return "Error en el proceso de inscripción"
}

// cleanErrorMessage limpia caracteres extra del mensaje de error
func cleanErrorMessage(message string) string {
	// Remover punto al inicio si existe
	message = strings.TrimPrefix(message, ".")
	// Remover espacios extra al inicio y final
	message = strings.TrimSpace(message)
	// Remover múltiples espacios
	message = strings.Join(strings.Fields(message), " ")

	// Agregar espacio entre "opcion" y números si no existe
	message = strings.ReplaceAll(message, "opcion", "opción")
	message = strings.ReplaceAll(message, "opción", "opción ")

	// Limpiar espacios múltiples nuevamente después de los reemplazos
	message = strings.Join(strings.Fields(message), " ")

	return message
}

// GetUserPassword obtiene la contraseña del usuario desde la base de datos
func (esp EnrollmentStoredProcedureDataSource) GetUserPassword(ctx context.Context, documentType, documentNumber string) (string, error) {
	query := `SELECT USR_CLAVE FROM COMUN.USUARIO 
			  WHERE TIPO_DOCUMENTO = :1 AND NUM_DOC_IDENTIDAD = :2`

	var password string
	err := esp.connection.QueryRowContext(ctx, query, documentType, documentNumber).Scan(&password) // TODO: Revisar si se puede usar el bind
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", &internalErrors.AppError{
				Message: "Usuario no encontrado en la base de datos",
				Err:     err,
			}
		}
		return "", cockroachdbErrors.Wrap(err, "infra: error obteniendo contraseña del usuario")
	}

	return password, nil
}

// verifyUserExists verifica si el usuario existe en la base de datos
//
// :param ctx: contexto de la aplicación
// :param conn: conexión a la base de datos
// :param nis: número de identificación del usuario
// :param documentType: tipo de documento
// :param documentNumber: número de documento
// :return: true si el usuario existe, false en caso contrario
func (esp EnrollmentStoredProcedureDataSource) verifyUserExists(ctx context.Context, conn *sql.Conn, nis int, documentType, documentNumber string) (bool, error) {
	query := `SELECT COUNT(*) FROM comun.usuario WHERE nis = :1 AND tipo_documento = :2 AND num_doc_identidad = :3`

	var count int
	err := conn.QueryRowContext(ctx, query, nis, documentType, documentNumber).Scan(&count)
	if err != nil {
		return false, cockroachdbErrors.Wrap(err, "infra: error verificando existencia del usuario")
	}

	return count > 0, nil
}

// verifyEnrollmentExists verifica si la inscripción existe en la base de datos
//
// :param ctx: contexto de la aplicación
// :param conn: conexión a la base de datos
// :param nis: número de identificación del usuario
// :param ficId: ID de la ficha de caracterización
// :return: true si la inscripción existe, false en caso contrario
func (esp EnrollmentStoredProcedureDataSource) verifyEnrollmentExists(ctx context.Context, conn *sql.Conn, nis int, ficId int) (bool, error) {
	query := `SELECT COUNT(*) FROM inscripcion.ingreso_aspirante WHERE nis = :1 AND fic_id = :2`

	var count int
	err := conn.QueryRowContext(ctx, query, nis, ficId).Scan(&count)
	if err != nil {
		return false, cockroachdbErrors.Wrap(err, "infra: error verificando existencia de la inscripción")
	}

	return count > 0, nil
}
