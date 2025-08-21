package postgresql

import (
	"context"
	"errors"
	"fmt"

	cockroachdbErrors "github.com/cockroachdb/errors"
	errorDbConnector "gitlab.com/sofia-plus/go_db_connectors/errors"
	postgresqlPort "gitlab.com/sofia-plus/go_db_connectors/ports/repositories/postgresql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/constants"
	postgresqlModels "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/postgresql"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
	"gorm.io/gorm"
)

type BasicDataUserDataSource struct {
	connection *gorm.DB
	db         postgresqlPort.PostgresqlRepository
}

func NewBasicDataUserDataSource(db postgresqlPort.PostgresqlRepository, connection *gorm.DB) BasicDataUserDataSource {
	return BasicDataUserDataSource{
		connection: connection,
		db:         db,
	}
}

func (b BasicDataUserDataSource) GetBasicDataUserByNis(ctx context.Context,nis int) (postgresqlModels.BasicDataUserModel,error) {
	basicDataUserModel := postgresqlModels.BasicDataUserModel{}
	where := b.db.Where(`"NIS_PG"`, "=", nis)
	if err := b.db.Select(&basicDataUserModel, where); err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInPostgres) {
			return basicDataUserModel, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,"infra: datos basicos no encontrados en postgresql: %s",err.Error() )
		}
		return basicDataUserModel, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo los datos basicos")
	}
	fmt.Println("linea 36")
	return basicDataUserModel, nil
}

func (b BasicDataUserDataSource) UpdateSyncStatusWithID(nisFromOra,nisFromPos int) error {
	result := b.connection.Exec(`UPDATE common."DATOS_BASICOS_USUARIO" SET "state" = $1,"NIS" = $2 WHERE "NIS_PG"= $3`,constants.SyncedStatus,nisFromOra,nisFromPos)
	if result.Error != nil{
		return cockroachdbErrors.Wrap(result.Error, "infra: ocurrio un error modificando el estado y clave original del registro en datos basicos")
	}
	return nil
}
 