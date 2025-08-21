package postgresql

// ModalityFilterModel representa el modelo para filtros de modalidad
type ModalityFilterModel struct {
	FormationModality string `gorm:"column:formation_modality"`
	ProgramModality   string `gorm:"column:program_modality"`
	Count             int64  `gorm:"column:count"`
}

func (ModalityFilterModel) TableName() string {
	return "programs.v_available_programs"
}

// DurationFilterModel representa el modelo para filtros de duración
type DurationFilterModel struct {
	ID    int    `gorm:"column:id"`
	Text  string `gorm:"column:text"`
	Count int64  `gorm:"column:count"`
}

func (DurationFilterModel) TableName() string {
	return "programs.v_available_programs"
}

// LevelFilterModel representa el modelo para filtros de nivel
type LevelFilterModel struct {
	ID    int    `gorm:"column:formation_level_id"`
	Text  string `gorm:"column:formation_level"`
	Count int64  `gorm:"column:count"`
}

func (LevelFilterModel) TableName() string {
	return "programs.v_available_programs"
}

// TechnologyLineFilterModel representa el modelo para filtros de línea tecnológica
type TechnologyLineFilterModel struct {
	ID       int    `gorm:"column:id"`
	Text     string `gorm:"column:text"`
	Quantity int64  `gorm:"column:quantity"`
}

func (TechnologyLineFilterModel) TableName() string {
	return "programs.v_available_programs"
}
