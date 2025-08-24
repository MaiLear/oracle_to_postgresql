package usecases

import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	repositoryPort "gitlab.com/sofia-plus/oracle_to_postgresql/domain/ports/repositories"
)

type TrainingProgram struct{
	repositoryReader repositoryPort.RepositoryReader
	repositoryWriter repositoryPort.RepositoryWrite
}

func (t TrainingProgram) SynchronizeData(ctx context.Context)(err error){
	programs,err := t.repositoryReader.Get(ctx)
	if err != nil{
		err = cockroachdbErrors.WithStack(err)
		return
	}
	if err = t.repositoryWriter.Save(ctx,programs); err != nil{
		err = cockroachdbErrors.WithStack(err)
		return
	}
	return
}