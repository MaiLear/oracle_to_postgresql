package postgresql

import "time"

// ParameterModel representa la tabla PARAMETRO en PostgreSQL
type ParameterModel struct {
	ParId          int        `gorm:"column:\"PAR_ID\""`
	ParNombre      string     `gorm:"column:\"PAR_NOMBRE\""`
	ParTipo        string     `gorm:"column:\"PAR_TIPO\""`
	ParDescripcion string     `gorm:"column:\"PAR_DESCRIPCION\""`
	ParValor       string     `gorm:"column:\"PAR_VALOR\""`
	PaqId          int        `gorm:"column:\"PAQ_ID\""`
	ParFchRegistro time.Time  `gorm:"column:\"PAR_FCH_REGISTRO\""`
	ParFchModif    *time.Time `gorm:"column:\"PAR_FCH_MODIF\""`
}

func (p ParameterModel) TableName() string {
	return "common.PARAMETRO"
}
