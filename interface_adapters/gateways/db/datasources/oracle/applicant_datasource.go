package oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	cockroachdbErrors "github.com/cockroachdb/errors"
	oraclePort "gitlab.com/sofia-plus/go_db_connectors/ports/repositories/oracle"
	oracleModels "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/oracle"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/tools"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
	errorDbConnector "gitlab.com/sofia-plus/go_db_connectors/errors"
)

type ApplicantDataSource struct {
	connection *sql.DB
	db         oraclePort.OracleRepository
}

func NewApplicantDataSource(db oraclePort.OracleRepository, connection *sql.DB) ApplicantDataSource {
	return ApplicantDataSource{
		connection: connection,
		db:         db,
	}
}

func (a ApplicantDataSource) Create(applicant oracleModels.ApplicantModel) error {
	if err := a.db.Insert(&applicant, nil); err != nil {
		if tools.InOraUniqueViolation(err) {
			return cockroachdbErrors.Wrap(internalErrors.ErrViolationUnique,fmt.Sprintf("infra: el aspirante ya se encuentra registrado %s",err.Error()))
		}
		return cockroachdbErrors.Wrap(err, "infra: ocurrio un error registrando el aspirante")
	}
	return nil
}

func (a ApplicantDataSource) GetApplicantByNis(ctx context.Context, nis int) (applicant *oracleModels.ApplicantModel, err error) {
	applicant = &oracleModels.ApplicantModel{}
	where := a.db.Where("NIS = :1", nis)
	_, err = a.db.Select(applicant, where)
	if err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInOracle) {
			return applicant, cockroachdbErrors.Wrap(internalErrors.ErrNotFound, fmt.Sprintf("infra: aspirante no encontrado en oracle %s", err.Error()))
		}
		return applicant, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo el aspirante en oracle")
	}
	return applicant, nil
}
