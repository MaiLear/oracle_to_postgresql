package oracle

import (
	"context"
	"database/sql"

	cockroachdbErrors "github.com/cockroachdb/errors"
	oracleDataSourcePort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/oracle"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
)

type EnrollmentValidationDataSource struct {
	connection *sql.DB
}

func NewEnrollmentValidationDataSource(connection *sql.DB) EnrollmentValidationDataSource {
	return EnrollmentValidationDataSource{connection: connection}
}

// GetUserByCredentials obtiene el NIS del usuario por credenciales
//
// :param ctx: contexto de la aplicación
// :param documentType: tipo de documento
// :param documentNumber: número de documento
// :param password: contraseña del usuario
// :return: NIS del usuario y error
func (evds EnrollmentValidationDataSource) GetUserByCredentials(ctx context.Context, documentType, documentNumber, password string) (int, error) {
	query := `SELECT NIS FROM COMUN.USUARIO 
			  WHERE TIPO_DOCUMENTO = :1 AND NUM_DOC_IDENTIDAD = :2 AND USR_CLAVE = :3`

	var nis int
	err := evds.connection.QueryRowContext(ctx, query, documentType, documentNumber, password).Scan(&nis)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, &internalErrors.AppError{
				Message: "La información de usuario y contraseña ingresada no es correcta.",
				Err:     err,
			}
		}
		return 0, cockroachdbErrors.Wrap(err, "infra: error obteniendo usuario por credenciales")
	}

	return nis, nil
}

// GetUserByNis obtiene el tipo de documento del usuario por NIS
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :return: tipo de documento y error
func (evds EnrollmentValidationDataSource) GetUserByNis(ctx context.Context, nis int) (string, error) {
	query := `SELECT TIPO_DOCUMENTO FROM COMUN.USUARIO WHERE NIS = :1`

	var documentType string
	err := evds.connection.QueryRowContext(ctx, query, nis).Scan(&documentType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &internalErrors.AppError{
				Message: "No se encontro información de usuario.",
				Err:     err,
			}
		}
		return "", cockroachdbErrors.Wrap(err, "infra: error obteniendo usuario por NIS")
	}

	return documentType, nil
}

// GetCourseInfo obtiene información del curso para validación
//
// :param ctx: contexto de la aplicación
// :param courseId: ID del curso
// :return: información del curso y error
func (evds EnrollmentValidationDataSource) GetCourseInfo(ctx context.Context, courseId int) (*oracleDataSourcePort.CourseValidationInfo, error) {
	query := `SELECT FC.FIC_ID, FC.PRF_ID, FC.FIC_LUGAR_REALIZACION, 
			  NVL(FC.FIC_PRUEBA_PREINSCRIPCION, 0), FC.NIP_ID
			  FROM PLANFORMACION.FICHA_CARACTERIZACION FC
			  WHERE FC.FIC_ID = :1`

	var courseInfo oracleDataSourcePort.CourseValidationInfo
	err := evds.connection.QueryRowContext(ctx, query, courseId).Scan(
		&courseInfo.CourseId,
		&courseInfo.ProgramId,
		&courseInfo.PlaceOfRealization,
		&courseInfo.TestPreEnrollment,
		&courseInfo.FormationLevelId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &internalErrors.AppError{
				Message: "La ficha de caracterización no existe.",
				Err:     err,
			}
		}
		return nil, cockroachdbErrors.Wrap(err, "infra: error obteniendo información del curso")
	}

	return &courseInfo, nil
}

// GetProgramInfo obtiene información del programa de formación
//
// :param ctx: contexto de la aplicación
// :param courseId: ID del curso
// :return: información del programa y error
func (evds EnrollmentValidationDataSource) GetProgramInfo(ctx context.Context, courseId int) (*oracleDataSourcePort.ProgramValidationInfo, error) {
	query := `SELECT PF.PRF_ID, PF.PRF_TIPO_PROGRAMA, PF.FLP_ID
			  FROM DISENIOCUR.PROGRAMA_FORMACION PF
			  INNER JOIN PLANFORMACION.FICHA_CARACTERIZACION FC ON FC.PRF_ID = PF.PRF_ID
			  WHERE FC.FIC_ID = :1`

	var programInfo oracleDataSourcePort.ProgramValidationInfo
	err := evds.connection.QueryRowContext(ctx, query, courseId).Scan(
		&programInfo.ProgramId,
		&programInfo.ProgramType,
		&programInfo.FamilyId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &internalErrors.AppError{
				Message: "La ficha de caracterización no existe.",
				Err:     err,
			}
		}
		return nil, cockroachdbErrors.Wrap(err, "infra: error obteniendo información del programa")
	}

	return &programInfo, nil
}

// GetNationalCourseForCompany obtiene el curso nacional correspondiente a un curso empresarial
//
// :param ctx: contexto de la aplicación
// :param companyCourseId: ID del curso empresarial
// :return: ID del curso nacional y error
func (evds EnrollmentValidationDataSource) GetNationalCourseForCompany(ctx context.Context, companyCourseId int) (int, error) {
	query := `SELECT DISTINCT FC.FIC_ID
			  FROM PLANFORMACION.FICHA_CARACTERIZACION FC
			  INNER JOIN PLANFORMACION.FICHA_CARACTERIZACION FCEMPRESA ON FCEMPRESA.PRF_ID = FC.PRF_ID
			  WHERE FCEMPRESA.FIC_ID = :1 AND
			  FC.FIC_ESTADO = 4 AND
			  FC.FIC_LUGAR_REALIZACION = 'VN' AND
			  FCEMPRESA.FIC_LUGAR_REALIZACION = 'VE' AND
			  ROWNUM = 1`

	var nationalCourseId int
	err := evds.connection.QueryRowContext(ctx, query, companyCourseId).Scan(&nationalCourseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, &internalErrors.AppError{
				Message: "Señor Usuario, En este momento el programa al cual está tratando de inscribirse no se encuentra disponible, lo invitamos a consultar los demás programas publicados en la oferta",
				Err:     err,
			}
		}
		return 0, cockroachdbErrors.Wrap(err, "infra: error obteniendo curso nacional para empresa")
	}

	return nationalCourseId, nil
}

// GetIdentityConfrontationStatus obtiene el estado de confrontación de identidad
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :return: estado de verificación y error
func (evds EnrollmentValidationDataSource) GetIdentityConfrontationStatus(ctx context.Context, nis int) (*oracleDataSourcePort.IdentityConfrontationInfo, error) {
	// Primero verificar si está habilitada la confrontación
	var habilitarConfrontacion int
	queryParam := `SELECT PAR_VALOR FROM COMUN.PARAMETRO WHERE PAR_ID = 259`
	err := evds.connection.QueryRowContext(ctx, queryParam).Scan(&habilitarConfrontacion)
	if err != nil {
		if err == sql.ErrNoRows {
			habilitarConfrontacion = 0
		} else {
			return nil, cockroachdbErrors.Wrap(err, "infra: error obteniendo parámetro de confrontación")
		}
	}

	if habilitarConfrontacion != 1 {
		return nil, nil // No está habilitada la confrontación
	}

	query := `SELECT ACN_ESTADO_VERIFICACION
			  FROM SEGURIDAD.ANALISIS_CONFRONTACION
			  WHERE NIS = :1`

	var confrontationInfo oracleDataSourcePort.IdentityConfrontationInfo
	err = evds.connection.QueryRowContext(ctx, query, nis).Scan(&confrontationInfo.VerificationStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No hay registro de confrontación
		}
		return nil, cockroachdbErrors.Wrap(err, "infra: error obteniendo estado de confrontación")
	}

	return &confrontationInfo, nil
}

// IsUserInstructorInProgram verifica si el usuario es instructor en el programa
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param programId: ID del programa
// :return: true si es instructor, false en caso contrario, y error
func (evds EnrollmentValidationDataSource) IsUserInstructorInProgram(ctx context.Context, nis int, programId int) (bool, error) {
	query := `SELECT COUNT(1)
			  FROM GESAMBIENTE.INSTRUCTORXFICHA INF
			  INNER JOIN PLANFORMACION.FICHA_CARACTERIZACION FIC ON FIC.FIC_ID = INF.FIC_ID
			  INNER JOIN DISENIOCUR.PROGRAMA_FORMACION PRF ON PRF.PRF_ID = FIC.PRF_ID
			  WHERE INF.NIS_FUN_INSTRUCTOR = :1
			  AND INF.INF_ESTADO = 'V'
			  AND PRF.PRF_ID = :2`

	var count int
	err := evds.connection.QueryRowContext(ctx, query, nis, programId).Scan(&count)
	if err != nil {
		return false, cockroachdbErrors.Wrap(err, "infra: error verificando si usuario es instructor")
	}

	return count > 0, nil
}

// IsUserInstructorInFamilyProgram verifica si el usuario es instructor en un programa de familia
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param familyId: ID de la familia de programas
// :return: true si es instructor, false en caso contrario, y error
func (evds EnrollmentValidationDataSource) IsUserInstructorInFamilyProgram(ctx context.Context, nis int, familyId int) (bool, error) {
	query := `SELECT COUNT(1)
			  FROM GESAMBIENTE.INSTRUCTORXFICHA INF
			  INNER JOIN PLANFORMACION.FICHA_CARACTERIZACION FIC ON FIC.FIC_ID = INF.FIC_ID
			  INNER JOIN DISENIOCUR.PROGRAMA_FORMACION PRF ON PRF.PRF_ID = FIC.PRF_ID
			  INNER JOIN DISENIOCUR.FAMILIAS_PROGRAMA FLP ON FLP.FLP_ID = PRF.FLP_ID
			  WHERE INF.NIS_FUN_INSTRUCTOR = :1
			  AND INF.INF_ESTADO = 'V'
			  AND FLP.FLP_ID = :2`

	var count int
	err := evds.connection.QueryRowContext(ctx, query, nis, familyId).Scan(&count)
	if err != nil {
		return false, cockroachdbErrors.Wrap(err, "infra: error verificando si usuario es instructor en familia")
	}

	return count > 0, nil
}

// GetUserEnrollmentsByProgram obtiene las inscripciones del usuario por programa
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param programId: ID del programa
// :return: estados de inscripción y error
func (evds EnrollmentValidationDataSource) GetUserEnrollmentsByProgram(ctx context.Context, nis int, programId int) ([]int, error) {
	query := `SELECT DISTINCT IA.ING_ESTADO
			  FROM PLANFORMACION.FICHA_CARACTERIZACION FC,
			  INSCRIPCION.INGRESO_ASPIRANTE IA
			  WHERE IA.FIC_ID = FC.FIC_ID
			  AND FC.PRF_ID = :1
			  AND IA.NIS = :2`

	rows, err := evds.connection.QueryContext(ctx, query, programId, nis)
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "infra: error obteniendo inscripciones por programa")
	}
	defer rows.Close()

	var states []int
	for rows.Next() {
		var state int
		if err := rows.Scan(&state); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "infra: error escaneando estado de inscripción")
		}
		states = append(states, state)
	}

	return states, nil
}

// GetUserVirtualEnrollmentsCount obtiene el conteo de inscripciones virtuales del usuario
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :return: número de inscripciones virtuales y error
func (evds EnrollmentValidationDataSource) GetUserVirtualEnrollmentsCount(ctx context.Context, nis int) (int, error) {
	query := `SELECT COUNT(1)
			  FROM INSCRIPCION.INGRESO_ASPIRANTE IA,
			  PLANFORMACION.FICHA_CARACTERIZACION FC,
			  DISENIOCUR.PROGRAMA_FORMACION PF
			  WHERE IA.NIS = :1
			  AND IA.ING_ESTADO IN (1,7,6,26,27)
			  AND FC.FIC_ID = IA.FIC_ID
			  AND PF.PRF_ID = FC.PRF_ID
			  AND PF.NFS_ID_OFRECIDO = 11`

	var count int
	err := evds.connection.QueryRowContext(ctx, query, nis).Scan(&count)
	if err != nil {
		return 0, cockroachdbErrors.Wrap(err, "infra: error obteniendo conteo de inscripciones virtuales")
	}

	return count, nil
}

// GetUserEnrollmentsByCourse obtiene las inscripciones del usuario por curso
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param courseId: ID del curso
// :return: estados de inscripción y error
func (evds EnrollmentValidationDataSource) GetUserEnrollmentsByCourse(ctx context.Context, nis int, courseId int) ([]int, error) {
	query := `SELECT ING_ESTADO
			  FROM INSCRIPCION.INGRESO_ASPIRANTE
			  WHERE NIS = :1 AND FIC_ID = :2 AND ING_ESTADO IN (1,7,26,27)`

	rows, err := evds.connection.QueryContext(ctx, query, nis, courseId)
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "infra: error obteniendo inscripciones por curso")
	}
	defer rows.Close()

	var states []int
	for rows.Next() {
		var state int
		if err := rows.Scan(&state); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "infra: error escaneando estado de inscripción")
		}
		states = append(states, state)
	}

	return states, nil
}

// GetUserEnrollmentsByFamilyProgram obtiene las inscripciones del usuario por familia de programas
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param familyId: ID de la familia de programas
// :return: estados de inscripción y error
func (evds EnrollmentValidationDataSource) GetUserEnrollmentsByFamilyProgram(ctx context.Context, nis int, familyId int) ([]int, error) {
	query := `SELECT COUNT(1)
			  FROM DISENIOCUR.PROGRAMA_FORMACION PF
			  INNER JOIN PLANFORMACION.FICHA_CARACTERIZACION FC ON FC.PRF_ID = PF.PRF_ID
			  INNER JOIN INSCRIPCION.INGRESO_ASPIRANTE IA ON IA.FIC_ID = FC.FIC_ID
			  WHERE PF.FLP_ID = :1 AND
			  IA.NIS = :2 AND
			  IA.ING_ESTADO IN (1,7,26,27,6)`

	var count int
	err := evds.connection.QueryRowContext(ctx, query, familyId, nis).Scan(&count)
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "infra: error obteniendo inscripciones por familia")
	}

	if count > 0 {
		return []int{1}, nil // Retornar un estado activo si hay inscripciones
	}
	return []int{}, nil
}

// GetAcademicRecordsByProgram obtiene los registros académicos del usuario por programa
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param programId: ID del programa
// :return: registros académicos y error
func (evds EnrollmentValidationDataSource) GetAcademicRecordsByProgram(ctx context.Context, nis int, programId int) ([]*oracleDataSourcePort.AcademicRecordInfo, error) {
	query := `SELECT RA.RGA_ID, RA.RGA_ESTADO,
			  DECODE(RA.RGA_ESTADO, 2, 'CANCELADO', 6, 'RETIRO VOLUNTARIO') ESTADO_NOMBRE,
			  RA.FIC_ID||'-'||PF.PRF_DENOMINACION NOMBRE_FICHA
			  FROM MATRICULA.REGISTRO_ACADEMICO RA,
			  DISENIOCUR.PROGRAMA_FORMACION PF
			  WHERE RA.NIS = :1
			  AND PF.PRF_ID = RA.PRF_ID
			  AND PF.PRF_TIPO_PROGRAMA = 'C'`

	rows, err := evds.connection.QueryContext(ctx, query, nis)
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "infra: error obteniendo registros académicos")
	}
	defer rows.Close()

	var records []*oracleDataSourcePort.AcademicRecordInfo
	for rows.Next() {
		var record oracleDataSourcePort.AcademicRecordInfo
		if err := rows.Scan(&record.RecordId, &record.Status, &record.StatusName, &record.ProgramName); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "infra: error escaneando registro académico")
		}
		records = append(records, &record)
	}

	return records, nil
}

// GetAcademicRecordNovelties obtiene las novedades de un registro académico
//
// :param ctx: contexto de la aplicación
// :param academicRecordId: ID del registro académico
// :return: novedades y error
func (evds EnrollmentValidationDataSource) GetAcademicRecordNovelties(ctx context.Context, academicRecordId int) ([]*oracleDataSourcePort.AcademicNoveltyInfo, error) {
	query := `SELECT SN.SUN_NOMBRE, NA.NAP_ESTADO, NA.NAP_FCH_ACTIVACION, 
			  NA.NAP_DURACION, SN.SUN_VALIDA_INSCRIPCION
			  FROM EJECFORMACION.NOVEDAD_APRENDIZ NA,
			  EJECFORMACION.SUBTIPO_NOVEDAD SN
			  WHERE NA.NAP_ID = (
			  SELECT MAX(NA.NAP_ID)
			  FROM EJECFORMACION.NOVEDAD_APRENDIZ NA
			  WHERE NA.RGA_ID = :1
			  AND NA.NAP_ESTADO IN (1,3,4))
			  AND SN.SUN_ID = NA.SUN_ID`

	rows, err := evds.connection.QueryContext(ctx, query, academicRecordId)
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "infra: error obteniendo novedades académicas")
	}
	defer rows.Close()

	var novelties []*oracleDataSourcePort.AcademicNoveltyInfo
	for rows.Next() {
		var novelty oracleDataSourcePort.AcademicNoveltyInfo
		if err := rows.Scan(&novelty.NoveltyTypeName, &novelty.Status, &novelty.ActivationDate, &novelty.Duration, &novelty.ValidatesEnrollment); err != nil {
			return nil, cockroachdbErrors.Wrap(err, "infra: error escaneando novedad académica")
		}
		novelties = append(novelties, &novelty)
	}

	return novelties, nil
}

// GetConfirmationDays obtiene los días de confirmación configurados
//
// :param ctx: contexto de la aplicación
// :return: días de confirmación y error
func (evds EnrollmentValidationDataSource) GetConfirmationDays(ctx context.Context) (string, error) {
	query := `SELECT PAR_VALOR FROM COMUN.PARAMETRO WHERE PAR_ID = 97`

	var confirmationDays string
	err := evds.connection.QueryRowContext(ctx, query).Scan(&confirmationDays)
	if err != nil {
		return "", cockroachdbErrors.Wrap(err, "infra: error obteniendo días de confirmación")
	}

	return confirmationDays, nil
}
