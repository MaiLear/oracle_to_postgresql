package postgresql


import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	dataSourcePort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/postgresql"
)

// BasicDataUserRepository implementa la interfaz de acceso a datos para Programas.
type BasicDataUserRepository struct {
	datasource dataSourcePort.BasicDataUserDataSource
}

func NewBasicDataUserRepository(datasource dataSourcePort.BasicDataUserDataSource) BasicDataUserRepository {
	return BasicDataUserRepository{datasource: datasource}
}

func (b BasicDataUserRepository) GetBasicDataUserByNis(ctx context.Context,nis int) (domainBasicDataUser entities.BasicUserData, err error) {
	basicDataUserModel, err := b.datasource.GetBasicDataUserByNis(ctx,nis)
	if err != nil {
		return domainBasicDataUser, cockroachdbErrors.WithStack(err)
	}
	domainBasicDataUser = mappers.FromBasicDataUserModelPosToDomain(basicDataUserModel)
	return domainBasicDataUser, nil
}

func (u BasicDataUserRepository) UpdateSyncStatusWithID(nisFromOra,nisFromPos int) error {
	if err := u.datasource.UpdateSyncStatusWithID(nisFromOra,nisFromPos); err != nil {
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}
