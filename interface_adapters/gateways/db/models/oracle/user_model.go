package oracle

import "time"

type UserModel struct {
	//Nis                      int        `db:"NIS"`
	DocumentType             string     `db:"TIPO_DOCUMENTO"`
	DocumentNumber           string     `db:"NUM_DOC_IDENTIDAD"`
	RegistrationDate         time.Time  `db:"USR_FCH_REGISTRO"`
	Password                 string     `db:"USR_CLAVE"`
	ExpirationDate           *time.Time `db:"USR_FCH_CADUCIDAD"`
	Branch                   *string    `db:"SUCURSAL"`
	Status                   int        `db:"USU_ESTADO"`
	AuthenticationDate       time.Time  `db:"USU_FCH_AUTENTICA"`
	IPAddress                string     `db:"USU_IP"`
	PasswordModDate          time.Time  `db:"USU_FCH_MOD_CLAV"`
	ActiveChangePassword     string     `db:"USU_ACTIVO_CAMBIAR_CLAVE"`
}
func (u UserModel) TableName() string {
	return "COMUN.USUARIO"
}
