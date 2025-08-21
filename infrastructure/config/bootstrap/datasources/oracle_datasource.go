package datasources

import (
	"database/sql"
	"log"

	datasource "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/oracle"
	dbConfig "gitlab.com/sofia-plus/go_db_connectors/config"
	connectorOracle "gitlab.com/sofia-plus/go_db_connectors/repositories/oracle"
)

type OracleDataSource struct {
	UserDataSource datasource.UserDataSource
	PeopleDataSource datasource.PeopleDataSource
	ApplicantDataSource datasource.ApplicantDataSource
	EnrollmentDataSource datasource.EnrollmentDataSource
	EnrollmentStoredProcedureDataSource datasource.EnrollmentStoredProcedureDataSource
	EnrollmentValidationDataSource datasource.EnrollmentValidationDataSource
	ParameterDataSource datasource.ParameterDataSource
	DocumentTypeDataSource datasource.DocumentTypeDatasource
}

type oracleConnector struct {
	UserConnector connectorOracle.Repository
	PeopleConnector connectorOracle.Repository
	ApplicantConnector connectorOracle.Repository
	EnrollmentConnector connectorOracle.Repository
	EnrollmentValidationConnector connectorOracle.Repository
	ParameterConector connectorOracle.Repository
}

func initOraConnectors() (oracleConnector, *sql.DB) {
	connection := initOraConnection()
	userConnector := connectorOracle.Repository{
		Connection: connection,
	}
	enrollmentConnector := connectorOracle.Repository{
		Connection: connection,
	}
	peopleConnector := connectorOracle.Repository{
		Connection: connection,
	}
	applicantConnector := connectorOracle.Repository{
		Connection: connection,
	}
	enrollmentValidationConnector := connectorOracle.Repository{
		Connection: connection,
	}
	parameterConnector := connectorOracle.Repository{
		Connection: connection,
	}

	connectors := oracleConnector{
		UserConnector: userConnector,
		PeopleConnector: peopleConnector,
		ApplicantConnector: applicantConnector,
		EnrollmentConnector: enrollmentConnector,
		EnrollmentValidationConnector: enrollmentValidationConnector,
		ParameterConector: parameterConnector,
	}

	return connectors, connection
}

func InitOraDataSorce() OracleDataSource {
	connectors, connection := initOraConnectors()

	userDataSource := datasource.NewUserDataSource(connectors.UserConnector, connection)
	enrollmentDataSource := datasource.NewEnrollmentDataSource(connectors.EnrollmentConnector, connection)
	enrollmentStoredProcedureDataSource := datasource.NewEnrollmentStoredProcedureDataSource(connection)
	peopleDataSource := datasource.NewPeopleDataSource(connectors.EnrollmentConnector, connection)
	applicantDataSource := datasource.NewApplicantDataSource(connectors.EnrollmentConnector, connection)
	enrollmentValidationDataSource := datasource.NewEnrollmentValidationDataSource(connection)
	parameterDataSource := datasource.NewParameterDataSource(connection)
	documentTypeDataSource := datasource.NewDocumentTypeDatasource(connection)
	return OracleDataSource{
		UserDataSource: userDataSource,
		PeopleDataSource: peopleDataSource,
		ApplicantDataSource: applicantDataSource,
		EnrollmentDataSource: enrollmentDataSource,
		EnrollmentStoredProcedureDataSource: enrollmentStoredProcedureDataSource,
		EnrollmentValidationDataSource: enrollmentValidationDataSource,
		ParameterDataSource: parameterDataSource,
		DocumentTypeDataSource: documentTypeDataSource,
	}
}

func initOraConnection() (connection *sql.DB) {
	// Crear conexión a base de datos
	connection, err := dbConfig.NewOracleConnection()
	if err != nil {
		log.Fatalf("no se pudo conectar a la bd postgresql %w", err)
	}
	return connection
}
