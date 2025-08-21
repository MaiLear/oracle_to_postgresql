package postgresql


import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	dataSourcePort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/postgresql"
)

// PeopleRepository implementa la interfaz de acceso a datos para Programas.
type PeopleRepository struct {
	datasource dataSourcePort.PeopleDataSource
}

func NewPeopleRepository(datasource dataSourcePort.PeopleDataSource) PeopleRepository {
	return PeopleRepository{datasource: datasource}
}

func (u PeopleRepository) GetPeopleByNis(ctx context.Context,nis int) (domainPeople entities.People, err error) {
	peopleModel, err := u.datasource.GetPeopleByNis(ctx,nis)
	if err != nil {
		return domainPeople, cockroachdbErrors.WithStack(err)
	}
	domainPeople = mappers.FromPeopleModelPosToDomain(peopleModel)
	return domainPeople, nil
}

func (p PeopleRepository) UpdateSyncStatusWithID(nisFromOra,nisFromPos int) error {
	if err := p.datasource.UpdateSyncStatusWithID(nisFromOra,nisFromPos); err != nil {
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}
