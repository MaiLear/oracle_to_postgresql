package usecases

import (
	"context"
	"errors"

	cockroachdbErrors "github.com/cockroachdb/errors"
	repositoryPort "gitlab.com/sofia-plus/oracle_to_postgresql/domain/ports/repositories"
)

type TrainingProgram struct {
	repositoryReader repositoryPort.RepositoryReader
	repositoryWriter repositoryPort.RepositoryWrite
}

func NewTrainingProgram(repositoryReader repositoryPort.RepositoryReader,repositoryWriter repositoryPort.RepositoryWrite)TrainingProgram{
	return TrainingProgram{
		repositoryReader: repositoryReader,
		repositoryWriter: repositoryWriter,
	}
}

func (t TrainingProgram) SynchronizeData(ctx context.Context) (err error) {
	var allErrors []error
	programs, err := t.repositoryReader.Get(ctx)
	if err != nil {
		err = cockroachdbErrors.WithStack(err)
		return
	}
	return
	for _,program := range programs{
		if ctx.Err() != nil{
			allErrors = append(allErrors, ctx.Err())
			break
		}
		if err = t.repositoryWriter.Upsert(ctx, program); err != nil {
			allErrors = append(allErrors, err)
		}
	}
	if len(allErrors)>0{
		err = cockroachdbErrors.WithStack(errors.Join(allErrors...))
	}
	return
}
