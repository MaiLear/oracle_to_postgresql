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

type UserLocationDataSource struct {
	connection *gorm.DB
	db         postgresqlPort.PostgresqlRepository
}

func NewUserLocationDataSource(db postgresqlPort.PostgresqlRepository, connection *gorm.DB) UserLocationDataSource {
	return UserLocationDataSource{
		connection: connection,
		db:         db,
	}
}

func (u UserLocationDataSource) GetUserLocationByNis(ctx context.Context,nis int) (postgresqlModels.UserLocationModel,error) {
	userLocationModel := postgresqlModels.UserLocationModel{}
	where := u.db.Where(`"NIS_PG"`, "=", nis)
	if err := u.db.Select(&userLocationModel, where); err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInPostgres) {
			return userLocationModel, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,"infra: ubicacion del usuario no encontrada en postgresql: %s",err.Error() )
		}
		return userLocationModel, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo la persona")
	}
	fmt.Println("linea 36")
	return userLocationModel, nil
}

func (u UserLocationDataSource) UpdateSyncStatusWithID(nisFromOra,nisFromPos int) error {
	result := u.connection.Exec(`UPDATE common."UBICACION_USUARIO" SET "state" = $1, "NIS" = $2 WHERE "NIS_PG" = $3`,constants.SyncedStatus,nisFromOra,nisFromPos)
	if result.Error != nil{
		return cockroachdbErrors.Wrap(result.Error, "infra: ocurrio un error modificando el estado y clave original del registro en user location")
	}
	return nil
}
 