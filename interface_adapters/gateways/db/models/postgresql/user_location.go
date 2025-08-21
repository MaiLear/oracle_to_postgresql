package postgresql

type UserLocationModel struct {
	NIS                int64  `gorm:"column:NIS;not null;primaryKey"`
	DocumentType       string `gorm:"column:TIPO_DOCUMENTO;size:3"`
	DocumentNumber     string `gorm:"column:NUM_DOC_IDENTIDAD;size:20"`
	Address            string `gorm:"column:UBU_DIR_RESIDENCIA;size:200"`
	CountryID          int64  `gorm:"column:PAI_ID_RESIDENCIA"`
	CountryName        string `gorm:"column:PAI_NOMBRE_RESIDENCIA;size:50"`
	DepartmentID       int64  `gorm:"column:DPT_ID_RESIDENCIA"`
	DepartmentName     string `gorm:"column:DPT_NOMBRE_RESIDENCIA;size:200"`
	MunicipalityID     int64  `gorm:"column:MPO_ID_RESIDENCIA"`
	MunicipalityName   string `gorm:"column:MPO_NOMBRE_RESIDENCIA;size:200"`
	ZoneID             int64  `gorm:"column:ZON_ID_RESIDENCIA"`
	ZoneName           string `gorm:"column:ZON_NOMBRE_RESIDENCIA;size:50"`
	NeighborhoodID     int64  `gorm:"column:BAR_ID_RESIDENCIA"`
	NeighborhoodName   string `gorm:"column:BAR_NOMBRE_RESIDENCIA;size:50"`
	MainPhone          string `gorm:"column:UBU_TEL_PRINCIPAL;size:20"`
	AlternativePhone   string `gorm:"column:UBU_TEL_ALTERNATIVO;size:20"`
	MobilePhone        string `gorm:"column:UBU_TEL_MOVIL;size:20"`
	EthnicGroup        int32  `gorm:"column:UBU_GRUPO_ETNICO"`
	LaborLinkage       int32  `gorm:"column:UBU_VINCULACION_LABORAL"`
	AddressInformation string `gorm:"column:UBU_INF_DIRECCION;size:100"`
}


func (UserLocationModel) TableName() string {
	return "common.UBICACION_USUARIO"
}
