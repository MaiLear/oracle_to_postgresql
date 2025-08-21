package postgresql

import "time"

type UserModel struct {
	Nis                  int        `gorm:"column:NIS;primaryKey"`
	NisLocal             int        `gorm:"column:NIS_PG"`
	DocumentType         string     `gorm:"column:TIPO_DOCUMENTO;not null"`
	DocumentNumber       string     `gorm:"column:NUM_DOC_IDENTIDAD;not null"`
	RegistrationDate     time.Time  `gorm:"column:USR_FCH_REGISTRO;not null"`
	Password             string     `gorm:"column:USR_CLAVE;not null"`
	ExpirationDate       *time.Time `gorm:"column:USR_FCH_CADUCIDAD"`
	Branch               *string    `gorm:"column:SUCURSAL"`
	Status               int        `gorm:"column:USU_ESTADO;default:2;not null"`
	AuthenticationDate   time.Time  `gorm:"column:USU_FCH_AUTENTICA;default:CURRENT_TIMESTAMP;not null"`
	IPAddress            string     `gorm:"column:USU_IP;default:'0.0.0.0';not null"`
	PasswordModDate      time.Time  `gorm:"column:USU_FCH_MOD_CLAV;default:CURRENT_TIMESTAMP;not null"`
	ActiveChangePassword string     `gorm:"column:USU_ACTIVO_CAMBIAR_CLAVE;default:'0';not null"`
	AttemptNumbers int 	`gorm:"column:number_attemps"`
}



func (UserModel) TableName()string{
	return "common.USUARIO"
}