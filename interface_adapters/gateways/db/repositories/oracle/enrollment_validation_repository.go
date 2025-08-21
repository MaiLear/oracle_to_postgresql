package oracle

import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	oracleDataSourcePort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/oracle"
	oracleDataSource "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/repositories/oracle"
)

type EnrollmentValidationRepository struct {
	datasource oracleDataSourcePort.EnrollmentValidationDataSource
}

func NewEnrollmentValidationRepository(datasource oracleDataSourcePort.EnrollmentValidationDataSource) EnrollmentValidationRepository {
	return EnrollmentValidationRepository{datasource: datasource}
}

// GetUserByCredentials obtiene el NIS del usuario por credenciales
//
// :param ctx: contexto de la aplicación
// :param documentType: tipo de documento
// :param documentNumber: número de documento
// :param password: contraseña del usuario
// :return: NIS del usuario y error
func (evr EnrollmentValidationRepository) GetUserByCredentials(ctx context.Context, documentType, documentNumber, password string) (int, error) {
	nis, err := evr.datasource.GetUserByCredentials(ctx, documentType, documentNumber, password)
	if err != nil {
		return 0, cockroachdbErrors.WithStack(err)
	}
	return nis, nil
}

// GetUserByNis obtiene el tipo de documento del usuario por NIS
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :return: tipo de documento y error
func (evr EnrollmentValidationRepository) GetUserByNis(ctx context.Context, nis int) (string, error) {
	documentType, err := evr.datasource.GetUserByNis(ctx, nis)
	if err != nil {
		return "", cockroachdbErrors.WithStack(err)
	}
	return documentType, nil
}

// GetCourseInfo obtiene información del curso para validación
//
// :param ctx: contexto de la aplicación
// :param courseId: ID del curso
// :return: información del curso y error
func (evr EnrollmentValidationRepository) GetCourseInfo(ctx context.Context, courseId int) (*oracleDataSource.CourseValidationInfo, error) {
	courseInfo, err := evr.datasource.GetCourseInfo(ctx, courseId)
	if err != nil {
		return nil, cockroachdbErrors.WithStack(err)
	}
	// Convertir de infraestructura a dominio
	domainCourseInfo := &oracleDataSource.CourseValidationInfo{
		CourseId:           courseInfo.CourseId,
		ProgramId:          courseInfo.ProgramId,
		PlaceOfRealization: courseInfo.PlaceOfRealization,
		TestPreEnrollment:  courseInfo.TestPreEnrollment,
		FormationLevelId:   courseInfo.FormationLevelId,
	}
	return domainCourseInfo, nil
}

// GetProgramInfo obtiene información del programa de formación
//
// :param ctx: contexto de la aplicación
// :param courseId: ID del curso
// :return: información del programa y error
func (evr EnrollmentValidationRepository) GetProgramInfo(ctx context.Context, courseId int) (*oracleDataSource.ProgramValidationInfo, error) {
	programInfo, err := evr.datasource.GetProgramInfo(ctx, courseId)
	if err != nil {
		return nil, cockroachdbErrors.WithStack(err)
	}
	// Convertir de infraestructura a dominio
	domainProgramInfo := &oracleDataSource.ProgramValidationInfo{
		ProgramId:   programInfo.ProgramId,
		ProgramType: programInfo.ProgramType,
		FamilyId:    programInfo.FamilyId,
	}
	return domainProgramInfo, nil
}

// GetNationalCourseForCompany obtiene el curso nacional correspondiente a un curso empresarial
//
// :param ctx: contexto de la aplicación
// :param companyCourseId: ID del curso empresarial
// :return: ID del curso nacional y error
func (evr EnrollmentValidationRepository) GetNationalCourseForCompany(ctx context.Context, companyCourseId int) (int, error) {
	nationalCourseId, err := evr.datasource.GetNationalCourseForCompany(ctx, companyCourseId)
	if err != nil {
		return 0, cockroachdbErrors.WithStack(err)
	}
	return nationalCourseId, nil
}

// GetIdentityConfrontationStatus obtiene el estado de confrontación de identidad
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :return: estado de verificación y error
func (evr EnrollmentValidationRepository) GetIdentityConfrontationStatus(ctx context.Context, nis int) (*oracleDataSource.IdentityConfrontationInfo, error) {
	confrontationInfo, err := evr.datasource.GetIdentityConfrontationStatus(ctx, nis)
	if err != nil {
		return nil, cockroachdbErrors.WithStack(err)
	}
	if confrontationInfo == nil {
		return nil, nil
	}
	// Convertir de infraestructura a dominio
	domainConfrontationInfo := &oracleDataSource.IdentityConfrontationInfo{
		VerificationStatus: confrontationInfo.VerificationStatus,
	}
	return domainConfrontationInfo, nil
}

// IsUserInstructorInProgram verifica si el usuario es instructor en el programa
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param programId: ID del programa
// :return: true si es instructor, false en caso contrario, y error
func (evr EnrollmentValidationRepository) IsUserInstructorInProgram(ctx context.Context, nis int, programId int) (bool, error) {
	isInstructor, err := evr.datasource.IsUserInstructorInProgram(ctx, nis, programId)
	if err != nil {
		return false, cockroachdbErrors.WithStack(err)
	}
	return isInstructor, nil
}

// IsUserInstructorInFamilyProgram verifica si el usuario es instructor en un programa de familia
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param familyId: ID de la familia de programas
// :return: true si es instructor, false en caso contrario, y error
func (evr EnrollmentValidationRepository) IsUserInstructorInFamilyProgram(ctx context.Context, nis int, familyId int) (bool, error) {
	isInstructor, err := evr.datasource.IsUserInstructorInFamilyProgram(ctx, nis, familyId)
	if err != nil {
		return false, cockroachdbErrors.WithStack(err)
	}
	return isInstructor, nil
}

// GetUserEnrollmentsByProgram obtiene las inscripciones del usuario por programa
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param programId: ID del programa
// :return: estados de inscripción y error
func (evr EnrollmentValidationRepository) GetUserEnrollmentsByProgram(ctx context.Context, nis int, programId int) ([]int, error) {
	states, err := evr.datasource.GetUserEnrollmentsByProgram(ctx, nis, programId)
	if err != nil {
		return nil, cockroachdbErrors.WithStack(err)
	}
	return states, nil
}

// GetUserVirtualEnrollmentsCount obtiene el conteo de inscripciones virtuales del usuario
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :return: número de inscripciones virtuales y error
func (evr EnrollmentValidationRepository) GetUserVirtualEnrollmentsCount(ctx context.Context, nis int) (int, error) {
	count, err := evr.datasource.GetUserVirtualEnrollmentsCount(ctx, nis)
	if err != nil {
		return 0, cockroachdbErrors.WithStack(err)
	}
	return count, nil
}

// GetUserEnrollmentsByCourse obtiene las inscripciones del usuario por curso
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param courseId: ID del curso
// :return: estados de inscripción y error
func (evr EnrollmentValidationRepository) GetUserEnrollmentsByCourse(ctx context.Context, nis int, courseId int) ([]int, error) {
	states, err := evr.datasource.GetUserEnrollmentsByCourse(ctx, nis, courseId)
	if err != nil {
		return nil, cockroachdbErrors.WithStack(err)
	}
	return states, nil
}

// GetUserEnrollmentsByFamilyProgram obtiene las inscripciones del usuario por familia de programas
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param familyId: ID de la familia de programas
// :return: estados de inscripción y error
func (evr EnrollmentValidationRepository) GetUserEnrollmentsByFamilyProgram(ctx context.Context, nis int, familyId int) ([]int, error) {
	states, err := evr.datasource.GetUserEnrollmentsByFamilyProgram(ctx, nis, familyId)
	if err != nil {
		return nil, cockroachdbErrors.WithStack(err)
	}
	return states, nil
}

// GetAcademicRecordsByProgram obtiene los registros académicos del usuario por programa
//
// :param ctx: contexto de la aplicación
// :param nis: número de identificación del usuario
// :param programId: ID del programa
// :return: registros académicos y error
func (evr EnrollmentValidationRepository) GetAcademicRecordsByProgram(ctx context.Context, nis int, programId int) ([]*oracleDataSource.AcademicRecordInfo, error) {
	records, err := evr.datasource.GetAcademicRecordsByProgram(ctx, nis, programId)
	if err != nil {
		return nil, cockroachdbErrors.WithStack(err)
	}
	// Convertir de infraestructura a dominio
	var domainRecords []*oracleDataSource.AcademicRecordInfo
	for _, record := range records {
		domainRecord := &oracleDataSource.AcademicRecordInfo{
			RecordId:    record.RecordId,
			Status:      record.Status,
			StatusName:  record.StatusName,
			CourseId:    record.CourseId,
			ProgramName: record.ProgramName,
		}
		domainRecords = append(domainRecords, domainRecord)
	}
	return domainRecords, nil
}

// GetAcademicRecordNovelties obtiene las novedades de un registro académico
//
// :param ctx: contexto de la aplicación
// :param academicRecordId: ID del registro académico
// :return: novedades y error
func (evr EnrollmentValidationRepository) GetAcademicRecordNovelties(ctx context.Context, academicRecordId int) ([]*oracleDataSource.AcademicNoveltyInfo, error) {
	novelties, err := evr.datasource.GetAcademicRecordNovelties(ctx, academicRecordId)
	if err != nil {
		return nil, cockroachdbErrors.WithStack(err)
	}
	// Convertir de infraestructura a dominio
	var domainNovelties []*oracleDataSource.AcademicNoveltyInfo
	for _, novelty := range novelties {
		domainNovelty := &oracleDataSource.AcademicNoveltyInfo{
			NoveltyTypeName:     novelty.NoveltyTypeName,
			Status:              novelty.Status,
			ActivationDate:      novelty.ActivationDate,
			Duration:            novelty.Duration,
			ValidatesEnrollment: novelty.ValidatesEnrollment,
		}
		domainNovelties = append(domainNovelties, domainNovelty)
	}
	return domainNovelties, nil
}

// GetConfirmationDays obtiene los días de confirmación configurados
//
// :param ctx: contexto de la aplicación
// :return: días de confirmación y error
func (evr EnrollmentValidationRepository) GetConfirmationDays(ctx context.Context) (string, error) {
	confirmationDays, err := evr.datasource.GetConfirmationDays(ctx)
	if err != nil {
		return "", cockroachdbErrors.WithStack(err)
	}
	return confirmationDays, nil
}
