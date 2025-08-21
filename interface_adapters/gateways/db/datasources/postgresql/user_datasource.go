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

type UserDataSource struct {
	connection *gorm.DB
	db         postgresqlPort.PostgresqlRepository
}

func NewUserDataSource(db postgresqlPort.PostgresqlRepository, connection *gorm.DB) UserDataSource {
	return UserDataSource{
		connection: connection,
		db:         db,
	}
}

func (u UserDataSource) GetPendingRecords(ctx context.Context) (usersWithPendingState []*postgresqlModels.UserModel, err error) {
	where := u.db.Where(`"state"`, "=", "pending")
	if err = u.db.SelectAll(&usersWithPendingState, where); err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInPostgres) {
			return usersWithPendingState, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,"infra: registros del usuario pendientes no encontrados en postgresql:  %s",err.Error() )
		}
		return usersWithPendingState, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo los registros pendientes del usuario")
	}
	return usersWithPendingState, nil
}

func (u UserDataSource) UpdateSyncStatusWithID(nisFromOra,nisFromPos int) error {
	result := u.connection.Exec(`UPDATE common."USUARIO" SET "state" = $1,"NIS" = $2 WHERE "NIS_PG"= $3`,constants.SyncedStatus,nisFromOra,nisFromPos)
	if result.Error != nil{
		return cockroachdbErrors.Wrap(result.Error, "infra: ocurrio un error modificando el estado y clave original del registro en usuario")
	}
	return nil
}

func (u UserDataSource) SetNumberAttemps(ctx context.Context, nisFromPos int) (err error) {
	result := u.connection.
		Model(&postgresqlModels.UserModel{}).
		Where(`"NIS_PG" = ?`, nisFromPos).
		Update("number_attemps",gorm.Expr("number_attemps + ?", 1))

	if result.Error != nil {
		return cockroachdbErrors.Wrapf(
			result.Error,
			"infra: ocurrió un error actualizando el numero de intentos con NIS_PG=%d: %s",
			nisFromPos, result.Error.Error(),
		)
	}
	if result.RowsAffected == 0 {
		return cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,
			"infra: no se encontró ningún registro con NIS_PG=%d para actualizar el numero de intentos",
			nisFromPos,
		)
	}

	return nil
}

func (u UserDataSource) UpdateEnrollmentStatusError(ctx context.Context, nisFromPos int) error {
	result := u.connection.
		Model(&postgresqlModels.UserModel{}).
		Where(`"NIS_PG" = ?`, nisFromPos).
		Update("state", "error")

	if result.Error != nil {
		return cockroachdbErrors.Wrapf(
			result.Error,
			"infra: ocurrió un error actualizando el status a 'error' para el ingreso con NIS_PG=%d: %s",
			nisFromPos, result.Error.Error(),
		)
	}
	if result.RowsAffected == 0 {
		return cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,
			"infra: no se encontró ningún registro con NIS_PG=%d para actualizar a 'error'",
			nisFromPos,
		)
	}

	return nil
}
