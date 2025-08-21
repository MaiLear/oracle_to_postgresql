package postgresql


type ErrorModel struct{
	Entity          string `gorm:"column:entity"`
	LocalID         int	`gorm:"column:local_id"`
	ErrorDescription string	`gorm:"column:error_description"`
	State           string  `gorm:"column:state"`
	ProcessName string  `gorm:"column:process_name"`
}

func (e ErrorModel) TableName() string {
	return "public.etl_error_log"
}


