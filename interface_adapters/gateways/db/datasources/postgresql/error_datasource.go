package postgresql

import (
	"context"
	cockroachdbErrors "github.com/cockroachdb/errors"
	postgresqlModels "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/postgresql"
	postgresqlPort "gitlab.com/sofia-plus/go_db_connectors/ports/repositories/postgresql"
	"gorm.io/gorm"
)

type ErrorDataSource struct {
	connection *gorm.DB
	db         postgresqlPort.PostgresqlRepository
}

func NewErrorDataSource(db postgresqlPort.PostgresqlRepository, connection *gorm.DB) ErrorDataSource {
	return ErrorDataSource{
		connection: connection,
		db:         db,
	}
}

func (e ErrorDataSource) SaveError(ctx context.Context,errorModel postgresqlModels.ErrorModel)error{
	if err := e.db.Insert(errorModel); err != nil{
		return cockroachdbErrors.Wrap(err,"infra: ocurrio un problema guardando el error")
	}
	return nil
}
