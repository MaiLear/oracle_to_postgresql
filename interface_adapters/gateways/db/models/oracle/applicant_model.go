package oracle

type ApplicantModel struct {
	Id               int    `db:"NIS"`
	DocumentType     string `db:"TIPO_DOCUMENTO"`
	NumDocument      string    `db:"NUM_DOC_IDENTIDAD"`
	Name             string `db:"ASP_NOMBRE"`
	FirstSurname     string `db:"ASP_PRIMER_APELLIDO"`
	SecondSurname    string `db:"ASP_SEGUNDO_APELLIDO"`
	Email            string `db:"ASP_CORREO_E"`
	AlternativeEmail string `db:"ASP_CORREO_ALT"`
	IdRegOfficer     *string `db:"NIS_FUN_REGISTRO"`
}

func (p ApplicantModel) TableName() string {
	return "COMUN.ASPIRANTE"
}
