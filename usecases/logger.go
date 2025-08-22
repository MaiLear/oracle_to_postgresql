package usecases

import (
	"context"

	"github.com/cockroachdb/errors"
	cockroachdbErrors "github.com/cockroachdb/errors"
	usecasesDto "gitlab.com/sofia-plus/oracle_to_postgresql/usecases/dto"
	logguerPort "gitlab.com/sofia-plus/oracle_to_postgresql/usecases/ports/out/loggers"
)

type Logger struct{
	logguers [] logguerPort.Logger
}

func (l Logger) Execute(ctx context.Context,dtoError usecasesDto.LogError)(err error){
	var allErrors []error
	for _,logger := range l.logguers{
		if err := logger.Save(ctx,dtoError); err != nil{
			allErrors = append(allErrors, err)
		}
	}
	if len(allErrors) > 0{
		err = cockroachdbErrors.WithStack(errors.Join(allErrors...))
		return
	}
	return
}