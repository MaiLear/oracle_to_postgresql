package oracle

type UserLocationModel struct {
	Nis                    int     `db:"NIS"`
	DocumentType           string  `db:"TIPO_DOCUMENTO"`
	DocumentNumber         string  `db:"NUM_DOC_IDENTIDAD"`
	Address                string  `db:"UBU_DIR_RESIDENCIA"`
	CountryID              int64     `db:"PAI_ID_RESIDENCIA"`
	CountryName            string  `db:"PAI_NOMBRE_RESIDENCIA"`
	DepartmentID           int64     `db:"DPT_ID_RESIDENCIA"`
	DepartmentName         string  `db:"DPT_NOMBRE_RESIDENCIA"`
	MunicipalityID         int64     `db:"MPO_ID_RESIDENCIA"`
	MunicipalityName       string  `db:"MPO_NOMBRE_RESIDENCIA"`
	ZoneID                 int64     `db:"ZON_ID_RESIDENCIA"`
	ZoneName               string  `db:"ZON_NOMBRE_RESIDENCIA"`
	NeighborhoodID         int64     `db:"BAR_ID_RESIDENCIA"`
	NeighborhoodName       string  `db:"BAR_NOMBRE_RESIDENCIA"`
	MainPhone              string  `db:"UBU_TEL_PRINCIPAL"`
	AlternativePhone       string  `db:"UBU_TEL_ALTERNATIVO"`
	MobilePhone            string  `db:"UBU_TEL_MOVIL"`
	EthnicGroup            int32  `db:"UBU_GRUPO_ETNICO"`
	LaborLinkage           int32  `db:"UBU_VINCULACION_LABORAL"`
	AddressInformation     string  `db:"UBU_INF_DIRECCION"`
}

func (u UserLocationModel) TableName() string {
	return "COMUN.UBICACION_USUARIO"
}
