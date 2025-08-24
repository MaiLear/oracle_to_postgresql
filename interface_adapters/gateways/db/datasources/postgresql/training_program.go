package postgresql

import (
	"context"
	"database/sql"

	oracleDbConnectorPort "gitlab.com/sofia-plus/go_db_connectors/ports/repositories/oracle"
	oracleModels "gitlab.com/sofia-plus/oracle_to_postgresql/interface_adapters/gateways/db/models/oracle"
)


type TrainingProgram struct{
	dbConnector oracleDbConnectorPort.OracleRepository
	connection *sql.DB
}

func (t TrainingProgram) updatePrograms(ctx context.Context,programs []oracleModels.TrainingProgram)(err error){
	
}

