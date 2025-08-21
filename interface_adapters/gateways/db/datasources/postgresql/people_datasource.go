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

type PeopleDataSource struct {
	connection *gorm.DB
	db         postgresqlPort.PostgresqlRepository
}

func NewPeopleDataSource(db postgresqlPort.PostgresqlRepository, connection *gorm.DB) PeopleDataSource {
	return PeopleDataSource{
		connection: connection,
		db:         db,
	}
}

func (p PeopleDataSource) GetPeopleByNis(ctx context.Context,nis int) (postgresqlModels.PeopleModel,error) {
	peopleModel := postgresqlModels.PeopleModel{}
	where := p.db.Where(`"NIS_PG"`, "=", nis)
	if err := p.db.Select(&peopleModel, where); err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInPostgres) {
			return peopleModel, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,"infra: persona no encontrada en postgresql:  %s",err.Error())
		}
		return peopleModel, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo la persona")
	}
	fmt.Println("linea 36")
	return peopleModel, nil
}

func (p PeopleDataSource) UpdateSyncStatusWithID(nisFromOra,nisFromPos int) error {
	result := p.connection.Exec(`UPDATE common."PERSONA" SET "state" = $1, "NIS" = $2 WHERE "NIS_PG"= $3`,constants.SyncedStatus,nisFromOra,nisFromPos)
	if result.Error != nil{
		return cockroachdbErrors.Wrap(result.Error, "infra: ocurrio un error modificando el estado y clave original del registro")
	}
	return nil
}
 