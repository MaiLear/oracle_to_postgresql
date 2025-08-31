package oracle

import (
	"context"
	//"database/sql"
	"fmt"

	"gitlab.com/sofia-plus/oracle_to_postgresql/domain/entities"
	"gorm.io/gorm"
)

type ProgramFormacion struct {
	PRFID   int64  `gorm:"column:PRF_ID;primaryKey"`
	Codigo  string `gorm:"column:PRF_CODIGO"`

	// Relación (has many)
	Fichas []FichaCaracterizacion `gorm:"foreignKey:PRF_ID;references:PRFID"`
}

type FichaCaracterizacion struct {
	ID     int64  `gorm:"column:FIC_ID;primaryKey"`
	PRFID  int64  `gorm:"column:PRF_ID"` // clave foránea
	Responsable string `gorm:"column:FIC_RESPONSABLE"`
}

func (f FichaCaracterizacion) TableName()string{
    return "PLANFORMACION.FICHA_CARACTERIZACION"
}

func (p ProgramFormacion) TableName()string{
    return "DISENIOCUR.PROGRAMA_FORMACION"
}


type TrainingProgram struct {
	dbConnection *gorm.DB
}

func NewTrainingProgram(dbConnection *gorm.DB) TrainingProgram {
	return TrainingProgram{
		dbConnection: dbConnection,
	}
}

// func (t TrainingProgram) Get(ctx context.Context) ([]entities.TrainingProgram, error) {
// 		var results []ProgramFormacion

// 	err := t.dbConnection.
// 		WithContext(ctx).
// 		Table(`"DISENIOCUR"."PROGRAMA_FORMACION"`).
// 		Preload("Fichas"). // 👈 aquí carga las fichas relacionadas
// 		Find(&results).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, r := range results {
// 		fmt.Printf("Programa: %+v\n", r)
// 		for _, f := range r.Fichas {
// 			fmt.Printf("  Ficha: %+v\n", f)
// 		}
// 	}

// 	return nil, nil
// }


func (t TrainingProgram) Get(ctx context.Context) ([]entities.TrainingProgram, error) { 
//     sqlDB, _ := t.dbConnection.DB()
// var count int
// _ = sqlDB.QueryRow(`SELECT COUNT(*) FROM "DISENIOCUR"."PROGRAMA_FORMACION"`).Scan(&count)
// fmt.Println("Filas:", count)
var results []ProgramFormacion 
fmt.Printf("DB connection: %#v\n", t.dbConnection) 
err := t.dbConnection. WithContext(ctx).Model(&ProgramFormacion{}).
Find(&results).Error
 if err != nil { 
 return nil, err 
 } 
 fmt.Printf("Resultados: %+v\n", results) 
 return nil, nil 
 }