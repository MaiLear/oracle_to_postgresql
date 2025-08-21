package postgresql

import (
	"time"
)

type CourseModel struct {
	Id          int       `gorm:"column:course_id"`
	Site        string    `gorm:"column:training_center"`
	City        string    `gorm:"column:town"`
	Department  string    `gorm:"column:department"`
	Modality    string    `gorm:"column:course_modality"`
	Shift       string    `gorm:"column:schedule_name"`
	StartDate   time.Time `gorm:"column:start_date"`
	ProgramId   int       `gorm:"column:program_id"`
	CourseState int       `gorm:"column:course_state"`
}

var allowedFields = []string{}
var allowedFilterFields = []string{
	"id", "site", "city", "department", "modality", "shift", "startdate", "programid",
}
var allowedOperators = []string{
	"=", "<>", "not in", "in",
}
var allowedExpand = []string{}
var allowedSort = []string{}
var allowedFieldSelections = []string{
	"*", "id", "site", "city", "department", "modality", "shift", "startdate",
}

func (c CourseModel) GetAllowedFilterFields() []string {
	return allowedFilterFields
}

func (c CourseModel) GetAllowedFieldSelection() []string {
	return allowedFieldSelections
}

func (c CourseModel) GetAllowedExpand() []string {
	return allowedExpand
}

func (c CourseModel) GetAllowedOperator() []string {
	return allowedOperators
}

func (CourseModel) TableName() string {
	return "programs.mv_available_courses"
}
