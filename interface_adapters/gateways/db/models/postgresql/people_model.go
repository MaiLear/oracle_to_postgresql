package postgresql

// PersonModel representa la tabla common.PERSONA en PostgreSQL
type PeopleModel struct {
	Nis              int     `gorm:"column:NIS;primaryKey"`
	DocumentType     string  `gorm:"column:TIPO_DOCUMENTO"`
	DocumentNumber   string  `gorm:"column:NUM_DOC_IDENTIDAD"`
	Name             string  `gorm:"column:PER_NOMBRE"`
	FirstSurname     string  `gorm:"column:PER_PRIMER_APELLIDO"`
	SecondSurname    *string `gorm:"column:PER_SEGUNDO_APELLIDO"`
	Email            string  `gorm:"column:PER_CORREO_E"`
	AlternativeEmail *string `gorm:"column:PER_CORREO_ALT"`
}

func (PeopleModel) TableName() string {
	return "common.PERSONA"
}