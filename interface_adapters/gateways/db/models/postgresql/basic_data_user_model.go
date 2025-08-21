package postgresql

import "time"

type BasicDataUserModel struct {
	NIS                     int64     `gorm:"column:NIS;not null"`
	DocumentType            string    `gorm:"column:TIPO_DOCUMENTO;size:3;not null"`
	DocumentNumber          string    `gorm:"column:NUM_DOC_IDENTIDAD;size:20;not null"`
	DocumentIssueDate       *time.Time `gorm:"column:DBU_FCH_EXD_DOC_IDENTIDAD"`
	Gender                  string    `gorm:"column:DBU_GENERO;size:1;not null"`
	BirthDate               time.Time `gorm:"column:DBU_FCH_NACIMIENTO;not null"`
	BirthMunicipalityID     *int64    `gorm:"column:MPO_ID_NACIMIENTO"`
	BirthMunicipalityName   string    `gorm:"column:MPO_NOMBRE_NACIMIENTO;size:200"`
	MilitaryCard            int8      `gorm:"column:DBU_LIB_MILITAR;not null"`
	MaritalStatus           string    `gorm:"column:DBU_ESTADO_CIVIL;size:1;not null"`
	Stratum                 int8      `gorm:"column:DBU_ESTRATO;not null"`
	BloodType               string    `gorm:"column:DBU_TIPO_SANGRE;size:3;not null"`
	IsEpsAffiliate          string    `gorm:"column:DBU_AFILIADO_EPS;size:1"`
	EpsName                 string    `gorm:"column:DBU_EPS;size:100"`
	EmergencyContactName    string    `gorm:"column:DBU_NOMBRE_CONTACTO;size:200"`
	EmergencyPhoneLandline  string    `gorm:"column:DBU_TEL_FIJO_CONTACTO;size:20"`
	EmergencyPhoneMobile    string    `gorm:"column:DBU_TEL_MOVIL_CONTACTO;size:20"`
	EmergencyRelationship   string    `gorm:"column:DBU_PARENTESCO_CONTACTO;size:20"`
	EmergencyCompany        string    `gorm:"column:DBU_EMPRESA_CONTACTO;size:200"`
	EmergencyBirthDate      *time.Time `gorm:"column:DBU_FCH_NACIMIENTO_CONTACTO"`
	IcfesScore              float64   `gorm:"column:DBU_PUNTAJE_ICFES"`
	HasTechnicalHighSchool  string    `gorm:"column:DBU_ES_MEDIA_TECNICA;size:1"`
	DocumentExpMunicipalityID *int64   `gorm:"column:MPO_ID_EXP_DOC_IDENTIDAD"`
	ShareInfoWithSena       string    `gorm:"column:DBU_COMPARTIR_INFO_SENA;size:1;not null;default:1"`
	IsSisben                string    `gorm:"column:DBU_SISBEN;size:1"`
	SisbenScore             int8      `gorm:"column:DBU_SISBEN_CALIFICA"`
	ChildrenCount           int8      `gorm:"column:DBU_HIJOS"`
	Ess                     string    `gorm:"column:DBU_ESS;size:50"`
	MilitaryCardNumber      string    `gorm:"column:DBU_NRO_LIB_MILITAR;size:20"`
	ExpirationDate          *time.Time `gorm:"column:DBU_FCH_VENCIMIENTO"`
	NISPG                   int64     `gorm:"column:NIS_PG;primaryKey;autoIncrement"`
	State                   string    `gorm:"column:state;size:10;default:pending"`
  }
  

func (BasicDataUserModel) TableName() string {
	return "common.DATOS_BASICOS_USUARIO"
}
