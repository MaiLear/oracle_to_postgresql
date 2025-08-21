package oracle

import "time"

// ParameterModel representa la tabla COMUN.PARAMETRO en Oracle
type ParameterModel struct {
	ParId          int        `db:"PAR_ID"`
	ParNombre      string     `db:"PAR_NOMBRE"`
	ParTipo        string     `db:"PAR_TIPO"`
	ParDescripcion string     `db:"PAR_DESCRIPCION"`
	ParValor       string     `db:"PAR_VALOR"`
	PaqId          int        `db:"PAQ_ID"`
	ParFchRegistro time.Time  `db:"PAR_FCH_REGISTRO"`
	ParFchModif    *time.Time `db:"PAR_FCH_MODIF"`
}

func (p ParameterModel) TableName() string {
	return "COMUN.PARAMETRO"
}
