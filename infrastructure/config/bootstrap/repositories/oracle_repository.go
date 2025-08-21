package repositories

import (
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/config/bootstrap/datasources"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/repositories/oracle"
)

type OracleRepository struct {
	UserRepository         oracle.UserRepository
	PeopleRepository         oracle.PeopleRepository
	ApplicantRepository         oracle.ApplicantRepository
	EnrollmentRepository         oracle.EnrollmentRepository
	EnrollmentValidationRepository oracle.EnrollmentValidationRepository
	ParameterRepository oracle.ParameterRepository
	DocumentTypeRepository oracle.DocumentTypeRepository
}

func InitOraRepository() OracleRepository {
	dataSources := datasources.InitOraDataSorce()
	userRepository := oracle.NewUserRepository(dataSources.UserDataSource)
	enrollmentRepository := oracle.NewEnrollmentRepository(dataSources.EnrollmentDataSource,dataSources.EnrollmentStoredProcedureDataSource)
	peopleRepository := oracle.NewPeopleRepository(dataSources.PeopleDataSource)
	applicantRepository := oracle.NewApplicantRepository(dataSources.ApplicantDataSource)
	enrollmentValidationRepository := oracle.NewEnrollmentValidationRepository(dataSources.EnrollmentValidationDataSource)
	parameterRepository := oracle.NewParameterRepository(dataSources.ParameterDataSource)
	documentTypeRepository := oracle.NewDocumentTypeRepository(dataSources.DocumentTypeDataSource)
	return OracleRepository{
		UserRepository: userRepository,
		PeopleRepository: peopleRepository,
		ApplicantRepository: applicantRepository,
		EnrollmentRepository: enrollmentRepository,
		EnrollmentValidationRepository: enrollmentValidationRepository,
		ParameterRepository: parameterRepository,
		DocumentTypeRepository: documentTypeRepository,
	}
}
