package postgresql

import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	postgresDataSource "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/postgresql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
	domainDto "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/dto"
)

type ErrorRepository struct {
	datasource postgresDataSource.ErrorDataSource
}

func NewErrorRepository(datasource postgresDataSource.ErrorDataSource) ErrorRepository {
	return ErrorRepository{datasource: datasource}
}

func (c ErrorRepository) SaveError(ctx context.Context, errorDto domainDto.ErrorDto) error {
	errorModel := mappers.FromErroDtoToModel(errorDto)
	if err := c.datasource.SaveError(ctx,errorModel); err != nil{
		return cockroachdbErrors.Wrap(err,"infra: ocurrio un problema guardando el error")
	}
	return nil
}
