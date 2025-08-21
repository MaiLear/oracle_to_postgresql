package oracle

import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	oracleDataSource "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/oracle"

	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
)

type ApplicantRepository struct {
	datasource oracleDataSource.ApplicantDataSource
}

func NewApplicantRepository(datasource oracleDataSource.ApplicantDataSource) ApplicantRepository {
	return ApplicantRepository{datasource: datasource}
}

func (a ApplicantRepository) CreateApplicantWithDataPeople(ctx context.Context, people entities.People) error {
	applicantModel := mappers.FromOraPeopleDomainToApplicantModel(people)
	if err := a.datasource.Create(applicantModel); err != nil {
		return cockroachdbErrors.Wrap(err, "infra: ocurrion un error creando el aspirante")
	}
	return nil
}

func (a ApplicantRepository) GetApplicantByNis(ctx context.Context, nis int) (applicantDomain entities.Applicant, err error) {
	applicantModel, err := a.datasource.GetApplicantByNis(ctx, nis)
	if err != nil {
		return applicantDomain, cockroachdbErrors.WithStack(err)
	}
	applicantDomain = mappers.FromApplicantModelOraToDomain(applicantModel)
	return applicantDomain, nil
}
