package oracle

import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	oracleDataSource "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/oracle"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
)

type PeopleRepository struct {
	datasource oracleDataSource.PeopleDataSource
}

func NewPeopleRepository(datasource oracleDataSource.PeopleDataSource) PeopleRepository {
	return PeopleRepository{datasource: datasource}
}

func (p PeopleRepository) GetPeopleByNis(ctx context.Context, nis int) (peopleDomain entities.People, err error) {
	peopleModel, err := p.datasource.GetPeopleByNis(ctx, nis)
	if err != nil {
		return peopleDomain, cockroachdbErrors.WithStack(err)
	}
	peopleDomain = mappers.FromOraPeopleModelToDomain(peopleModel)
	return peopleDomain, nil

}
