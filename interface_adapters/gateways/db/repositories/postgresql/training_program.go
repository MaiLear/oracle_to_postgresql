package postgresql

import (
	"context"

	"gitlab.com/sofia-plus/oracle_to_postgresql/domain/entities"
	"gorm.io/gorm"
)

type TrainingProgram struct{
	dbConnection *gorm.DB
}

func NewTrainingProgram(dbConnection *gorm.DB) TrainingProgram{
	return TrainingProgram{
		dbConnection: dbConnection,
	}
}

func (t TrainingProgram) Upsert(context.Context, entities.TrainingProgram) error{
	return nil
}