package postgresql

type ProgramModel struct {
	ID                       int    `gorm:"column:program_id"`
	Code                     string `gorm:"column:program_code"`
	Name                     string `gorm:"column:program_name"`
	Description              string `gorm:"column:program_description"`
	Duration                 int    `gorm:"column:program_duration"`
	TechnologyLineID         int    `gorm:"column:tecnology_line_id"`
	TechnologyLine           string `gorm:"column:tecnology_line"`
	FormattedDuration        string `gorm:"column:formatted_duration"`
	ProgramType              string `gorm:"column:program_type"`
	ProgramTypeText          string `gorm:"column:program_type_text"`
	FormationModality        string `gorm:"column:formation_modality"`
	ModalityText             string `gorm:"column:program_modality"`
	FormationLevelID         int    `gorm:"column:formation_level_id"`
	FormationLevel           string `gorm:"column:formation_level"`
	Degree                   string `gorm:"column:degree"`
	DescriptionRequirements  string `gorm:"column:description_requirements"`
	AgeRequirement           int    `gorm:"column:age_requirement"`
	AcademicLevelRequirement string `gorm:"column:academic_level_requirement"`
	Competencies             string `gorm:"column:competencies"`
}

var allowedProgramFields = []string{}

var allowedProgramFilterFields = []string{
	"id",
	"code",
	"name",
	"duration",
	"programtype",
	"programtypetext",
	"formationmodality",
	"modalitytext",
	"formationlevelid",
	"formationlevel",
	"technologylineid",
	"technologyline",
	"degree",
}

var allowedProgramOperators = []string{
	"=", "<>", "not in", "in",
}

var allowedProgramExpand = []string{}

var allowedProgramSort = []string{}

var allowedProgramFieldSelections = []string{
	"*",
	"id",
	"code",
	"name",
	"description",
	"duration",
	"formattedduration",
	"programtype",
	"programtypetext",
	"formationmodality",
	"modalitytext",
	"formationlevelid",
	"formationlevel",
	"degree",
	"descriptionrequirements",
	"agerequirement",
	"academiclevelrequirement",
	"technologylineid",
	"technologyline",
	"competencies",
}

func (c ProgramModel) GetAllowedFilterFields() []string {
	return allowedProgramFilterFields
}

func (c ProgramModel) GetAllowedFieldSelection() []string {
	return allowedProgramFieldSelections
}

func (c ProgramModel) GetAllowedExpand() []string {
	return allowedProgramExpand
}

func (c ProgramModel) GetAllowedOperator() []string {
	return allowedProgramOperators
}

func (ProgramModel) TableName() string {
	return "programs.mv_available_programs"
}
