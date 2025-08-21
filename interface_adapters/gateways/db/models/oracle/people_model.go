package oracle

type PeopleModel struct {
	Id               int    `db:"NIS"`
	DocumentType     string `db:"TIPO_DOCUMENTO"`
	NumDocument      string    `db:"NUM_DOC_IDENTIDAD"`
	Name             string `db:"PER_NOMBRE"`
	FirstSurname     string `db:"PER_PRIMER_APELLIDO"`
	SecondSurname    string `db:"PER_SEGUNDO_APELLIDO"`
	Email            string `db:"PER_CORREO_E"`
	AlternativeEmail string `db:"PER_CORREO_ALT"`
}

func (p PeopleModel) TableName() string {
	return "COMUN.PERSONA"
}
