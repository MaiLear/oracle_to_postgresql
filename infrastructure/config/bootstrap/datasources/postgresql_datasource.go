package datasources

import (
	"log"

	datasource "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/postgresql"
	dbConfig "gitlab.com/sofia-plus/go_db_connectors/config"
	connectorPostgresql "gitlab.com/sofia-plus/go_db_connectors/repositories/postgresql"
	"gorm.io/gorm"
)

type postgreSqlDataSource struct {
	UserDataSource    datasource.UserDataSource
	PeopleDataSource datasource.PeopleDataSource
	BasicDataUserDataSource datasource.BasicDataUserDataSource
	UserLocationDataSource datasource.UserLocationDataSource
	EnrollmentDataSource datasource.EnrollmentDataSource
	ErrorDataSource   datasource.ErrorDataSource
	ApplicantDataSource   datasource.ApplicantDataSource
	ParameterDataSource datasource.ParameterDataSource
	DocumentTypeDataSource datasource.DocumentTypeDatasource
}

type postgreSqlConnector struct {
	UserConnector        connectorPostgresql.Repository
	PeopleConnector 	connectorPostgresql.Repository
	BasicDataUserConnector  connectorPostgresql.Repository
	UserLocationConnector connectorPostgresql.Repository
	EnrollmentConnector    connectorPostgresql.Repository
	ErrorConnector       connectorPostgresql.Repository
	ApplicantConnector       connectorPostgresql.Repository
}

func initPosConnectors() (postgreSqlConnector, *gorm.DB) {
	connection := initPosConnection()
	userConnector := connectorPostgresql.Repository{
		Connection: connection,
	}
	peopleConnector :=  connectorPostgresql.Repository{
		Connection: connection,
	}
	basicDataUserConnector :=  connectorPostgresql.Repository{
		Connection: connection,
	}
	userLocationConnector :=  connectorPostgresql.Repository{
		Connection: connection,
	}
	enrollmentConnector :=  connectorPostgresql.Repository{
		Connection: connection,
	}
	errorConnector :=  connectorPostgresql.Repository{
		Connection: connection,
	}
	applicantConnector :=  connectorPostgresql.Repository{
		Connection: connection,
	}

	return postgreSqlConnector{
		UserConnector: userConnector,
		PeopleConnector: peopleConnector,
		BasicDataUserConnector: basicDataUserConnector,
		UserLocationConnector: userLocationConnector,
		EnrollmentConnector: enrollmentConnector,
		ErrorConnector: errorConnector,
		ApplicantConnector: applicantConnector,
	}, connection
}

func InitPosDataSorce() postgreSqlDataSource {
	connectors, connection := initPosConnectors()
	userDataSource := datasource.NewUserDataSource(connectors.UserConnector, connection)
	peopleDataSource := datasource.NewPeopleDataSource(connectors.UserConnector, connection)
	basicDataUserDataSource := datasource.NewBasicDataUserDataSource(connectors.UserConnector, connection)
	userLocationDataSource := datasource.NewUserLocationDataSource(connectors.UserConnector, connection)
	enrollmentDataSource := datasource.NewEnrollmentDataSource(connectors.EnrollmentConnector, connection)
	errorDataSource := datasource.NewErrorDataSource(connectors.UserConnector, connection)
	applicantDataSource := datasource.NewApplicantDataSource(connectors.UserConnector, connection)
	parameterDataSource := datasource.NewParameterDataSource(connection)
	documentTypeDataSource := datasource.NewDocumentTypeDatasource(connection)
	return postgreSqlDataSource{
		UserDataSource: userDataSource,
		PeopleDataSource: peopleDataSource,
		BasicDataUserDataSource: basicDataUserDataSource,
		UserLocationDataSource: userLocationDataSource,
		EnrollmentDataSource: enrollmentDataSource,
		ErrorDataSource: errorDataSource,
		ApplicantDataSource: applicantDataSource,
		ParameterDataSource: parameterDataSource,
		DocumentTypeDataSource: documentTypeDataSource,
	}
}

func initPosConnection() (connection *gorm.DB) {
	connection, err := dbConfig.NewPostgresConnection()
	if err != nil {
		log.Fatalf("no se pudo conectar a la bd postgresql %w", err)
	}
	return connection
}
