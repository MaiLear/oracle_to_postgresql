package postgresql

type ApplicantModel struct {
	Id               int    `gorm:"column:NIS"`
	DocumentType     string `gorm:"column:TIPO_DOCUMENTO"`
	NumDocument      string `gorm:"column:NUM_DOC_IDENTIDAD"`
	Name             string `gorm:"column:ASP_NOMBRE"`
	FirstSurname     string `gorm:"column:ASP_PRIMER_APELLIDO"`
	SecondSurname    string `gorm:"column:ASP_SEGUNDO_APELLIDO"`
	Email            string `gorm:"column:ASP_CORREO_E"`
	AlternativeEmail string `gorm:"column:ASP_CORREO_ALT"`
	IdRegOfficer     *string `gorm:"column:NIS_FUN_REGISTRO"`
}

func (p ApplicantModel) TableName() string {
	return `common."ASPIRANTE"`
}