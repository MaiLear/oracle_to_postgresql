package oracle

import (
	"context"
	"database/sql"
	"errors"

	cockroachdbErrors "github.com/cockroachdb/errors"
	oracleModels "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/oracle"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
	errorDbConnector "gitlab.com/sofia-plus/go_db_connectors/errors"
	oraclePort "gitlab.com/sofia-plus/go_db_connectors/ports/repositories/oracle"
)

type PeopleDataSource struct {
	connection *sql.DB
	db         oraclePort.OracleRepository
}

func NewPeopleDataSource(db oraclePort.OracleRepository, connection *sql.DB) PeopleDataSource {
	return PeopleDataSource{
		connection: connection,
		db:         db,
	}
}

func (p PeopleDataSource) GetPeopleByNis(ctx context.Context, nis int) (people *oracleModels.PeopleModel, err error) {
	people = &oracleModels.PeopleModel{}
	where := p.db.Where("NIS = :1",nis)
	_, err = p.db.Select(people, where)
	if err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInOracle) {
			return people, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,"infra: no se encontro a la persona en oracle: %s",err.Error())
		}
		return people, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo la personas en datasource de oracle")
	}
	return people, nil
}
