package oracle



type ProgramModel struct {
	Id        int `db:"PRF_ID"`
	ProMinAge *int `db:"PRF_EDAD_MIN_REQUERIDA"`
	TypePro string `db:"PRF_TIPO_PROGRAMA"`
	TechnologyLineId int `db:"LTC_ID"`
}

func (ProgramModel) TableName() string {
	return "INTEGRACION.V_PROGRAMA_FORMACION_B PR"
}
