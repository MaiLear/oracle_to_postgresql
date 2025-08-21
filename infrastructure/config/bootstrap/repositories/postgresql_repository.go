package repositories

import (
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/config/bootstrap/datasources"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/repositories/postgresql"
)

type PostgresqlRepository struct {
	UserRepository        postgresql.UserRepository
	PeopleRepository        postgresql.PeopleRepository
	BasicDataUserRepository        postgresql.BasicDataUserRepository
	UserLocationRepository        postgresql.UserLocationRepository
	EnrollmentRepository 	  postgresql.EnrollmentRepository
	ErrorRepository       postgresql.ErrorRepository
	ApplicantRepository postgresql.ApplicantRepository
	ParameterRepository postgresql.ParameterRepository
	DocumentTypeRepository postgresql.DocumentTypeRepository
}

func InitPosRepository() PostgresqlRepository {
	dataSources := datasources.InitPosDataSorce()
	userRepository := postgresql.NewUserRepository(dataSources.UserDataSource)
	peopleRepository := postgresql.NewPeopleRepository(dataSources.PeopleDataSource)
	basicDataUserRepository := postgresql.NewBasicDataUserRepository(dataSources.BasicDataUserDataSource)
	userLocationRepository := postgresql.NewUserLocationRepository(dataSources.UserLocationDataSource)
	enrollmentRepository := postgresql.NewEnrollmentRepository(dataSources.EnrollmentDataSource)
	errorRepository := postgresql.NewErrorRepository(dataSources.ErrorDataSource)
	applicantRepository := postgresql.NewApplicantRepository(dataSources.ApplicantDataSource)
	parameterRepository := postgresql.NewParameterRepository(dataSources.ParameterDataSource)
	documentTypeRepository := postgresql.NewDocumentTypeRepository(dataSources.DocumentTypeDataSource)
	return PostgresqlRepository{
		UserRepository: userRepository,
		PeopleRepository: peopleRepository,
		BasicDataUserRepository: basicDataUserRepository,
		UserLocationRepository: userLocationRepository,
		EnrollmentRepository: enrollmentRepository,
		ErrorRepository: errorRepository,
		ApplicantRepository: applicantRepository,
		ParameterRepository: parameterRepository,
		DocumentTypeRepository: documentTypeRepository,
	}
}
