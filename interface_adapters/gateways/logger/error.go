package logger

import (
	"context"
	"fmt"

	domainDto "gitlab.com/sofia-plus/oracle_to_postgresql/domain/dto"
	postgresqlRepPort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/repositories/postgresql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/logger"
	cockroachdbErrors "github.com/cockroachdb/errors"
)

type Error struct {
	postgresqlRepo postgresqlRepPort.ErrorRepository
}

func NewError(postgresqlRepo postgresqlRepPort.ErrorRepository) Error {
	return Error{
		postgresqlRepo: postgresqlRepo,
	}
}

func (e Error) LogErrorInFile(errorDto domainDto.ErrorDto){
	message := fmt.Sprintf(
		"[ERROR] entity=%s, local_id=%d, description=%s, state=%s",
		errorDto.Entity,
		errorDto.LocalID,
		fmt.Sprintf("%+v", errorDto.Error),
		errorDto.State,
	)
	logger.ErrorLogger.Println(message)
}



func (e Error) ErrorLog(ctx context.Context,allErrors []domainDto.ErroItem,entityName string,state string)error{
	for _,err := range allErrors{
		errorDto := domainDto.ErrorDto{
			Entity: entityName,
			LocalID: err.LocalID,
			Error: err.Err,
			State: state,
		}
		e.LogErrorInFile(errorDto)
		if err := e.SaveError(ctx,errorDto); err != nil{
			return cockroachdbErrors.WithStack(err)
		}
	}
	return nil
}
