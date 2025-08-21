package postgresql

type VanueEventModel struct{
	CharacterizacionSheetId int `gorm:"column:FIC_ID"`
}

func (v VanueEventModel) TableName()string{
	return `admcalendar.EVENTO_SEDE`
}