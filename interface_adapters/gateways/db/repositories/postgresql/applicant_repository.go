package postgresql

import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	dataSourcePort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/postgresql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
)

// ApplicantRepository implementa la interfaz de acceso a datos para Programas.
type ApplicantRepository struct {
	datasource dataSourcePort.ApplicantDataSource
}

func NewApplicantRepository(datasource dataSourcePort.ApplicantDataSource) ApplicantRepository {
	return ApplicantRepository{datasource: datasource}
}

func (a ApplicantRepository) GetApplicantByNis(ctx context.Context,nis int) (applicantDomain entities.Applicant, err error) {
	applicantModel, err := a.datasource.GetApplicantByNis(ctx,nis)
	if err != nil {
		return applicantDomain, cockroachdbErrors.WithStack(err)
	}
	applicantDomain = mappers.FromApplicantModelPosToDomain(applicantModel)
	return applicantDomain, nil
}

func (a ApplicantRepository) UpdateSyncStatusWithID(nisFromOra,nisFromPos int) error {
	if err := a.datasource.UpdateSyncStatusWithID(nisFromOra,nisFromPos); err != nil {
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}

