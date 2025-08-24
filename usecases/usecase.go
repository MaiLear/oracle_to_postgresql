package usecases

import (
	"context"
	"errors"
	"time"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/oracle_to_postgresql/usecases/ports/pipeline"
)



type UseCase struct{
	usecases []pipeline.Service
}


func (u UseCase) Execute()(err error){
	var allErrors []error
	ctx,cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()
	for _,usecase := range u.usecases{
		if err := usecase.SynchronizeData(ctx); err != nil{
			allErrors = append(allErrors, err)
		}
	}
	if len(allErrors) > 0{
		err = cockroachdbErrors.WithStack(errors.Join(allErrors...))
		return
	}
	return
}