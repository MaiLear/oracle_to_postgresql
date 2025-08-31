package usecases

import (
	"context"
	"errors"
	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/oracle_to_postgresql/usecases/ports/pipeline"
)

type UseCase struct {
	usecases []pipeline.Service
}

func NewUseCase(usecases []pipeline.Service)UseCase{
	return UseCase{
		usecases: usecases,
	}
}

func (u UseCase) Execute(ctx context.Context) (err error) {
	var allErrors []error
	for _, usecase := range u.usecases {
		if err := usecase.SynchronizeData(ctx); err != nil {
			allErrors = append(allErrors, err)
		}
	}
	if len(allErrors) > 0 {
		err = cockroachdbErrors.WithStack(errors.Join(allErrors...))
		return
	}
	return
}
