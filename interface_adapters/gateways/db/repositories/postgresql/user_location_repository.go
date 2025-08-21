package postgresql


import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	dataSourcePort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/postgresql"
)

// UserLocationRepository implementa la interfaz de acceso a datos para Programas.
type UserLocationRepository struct {
	datasource dataSourcePort.UserLocationDataSource
}

func NewUserLocationRepository(datasource dataSourcePort.UserLocationDataSource) UserLocationRepository {
	return UserLocationRepository{datasource: datasource}
}

func (u UserLocationRepository) GetUserLocationByNis(ctx context.Context,nis int) (domainUserLocation entities.UserLocation, err error) {
	userLocationModel, err := u.datasource.GetUserLocationByNis(ctx,nis)
	if err != nil {
		return domainUserLocation, cockroachdbErrors.WithStack(err)
	}
	domainUserLocation = mappers.FromUserLocationModelPosToDomain(userLocationModel)
	return domainUserLocation, nil
}

func (u UserLocationRepository) UpdateSyncStatusWithID(nisFromOra,nisFromPos int) error {
	if err := u.datasource.UpdateSyncStatusWithID(nisFromOra,nisFromPos); err != nil {
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}
