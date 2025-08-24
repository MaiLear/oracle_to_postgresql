package repositories

import (
	"context"

	"gitlab.com/sofia-plus/oracle_to_postgresql/domain/entities"
)

type RepositoryReader interface{
	Get(context.Context)([]entities.TrainingProgram,error)
}

type RepositoryWrite interface{
	Save(context.Context,[]entities.TrainingProgram)error
}