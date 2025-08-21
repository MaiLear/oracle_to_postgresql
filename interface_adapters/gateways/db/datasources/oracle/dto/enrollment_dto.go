package dto


type EnrollmentWithProgramAndCourseDto struct {
	CourseID int `gorm:"column:FIC_ID"`
	ProgramType string `gorm:"column:PROGRAM_TYPE"`
}

func (e EnrollmentWithProgramAndCourseDto) TableName() string {
	return "INSCRIPCION.INGRESO_ASPIRANTE"
}
