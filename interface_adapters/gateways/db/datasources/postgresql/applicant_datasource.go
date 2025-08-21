package postgresql

import (
	"context"
	"errors"
	cockroachdbErrors "github.com/cockroachdb/errors"
	errorDbConnector "gitlab.com/sofia-plus/go_db_connectors/errors"
	postgresqlPort "gitlab.com/sofia-plus/go_db_connectors/ports/repositories/postgresql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/constants"
	postgresqlModels "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/postgresql"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
	"gorm.io/gorm"
)

type ApplicantDataSource struct {
	connection *gorm.DB
	db         postgresqlPort.PostgresqlRepository
}

func NewApplicantDataSource(db postgresqlPort.PostgresqlRepository, connection *gorm.DB) ApplicantDataSource {
	return ApplicantDataSource{
		connection: connection,
		db:         db,
	}
}
func (a ApplicantDataSource) GetApplicantByNis(ctx context.Context,nis int) (postgresqlModels.ApplicantModel,error) {
	applicantModel := postgresqlModels.ApplicantModel{}
	where := a.db.Where(`"NIS_PG"`, "=", nis)
	if err := a.db.Select(&applicantModel, where); err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInPostgres) {
			return applicantModel, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,"infra: aspirante no encontrado en postgresql: %s",err.Error() )
		}
		return applicantModel, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo la persona")
	}
	return applicantModel, nil
}

func (a ApplicantDataSource) UpdateSyncStatusWithID(nisFromOra,nisFromPos int) error {
	result := a.connection.Exec(`UPDATE common."ASPIRANTE" SET "state" = $1,"NIS" = $2 WHERE "NIS_PG"= $3`,constants.SyncedStatus,nisFromOra,nisFromPos)
	if result.Error != nil{
		return cockroachdbErrors.Wrap(result.Error, "infra: ocurrio un error modificando el estado y clave original del registro en aspirante")
	}
	return nil
}

 