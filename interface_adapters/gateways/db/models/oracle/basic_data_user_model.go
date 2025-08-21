package oracle

import "time"

type BasicDataUserModel struct {
	Nis                      int      `db:"NIS"`
	DocumentType             string     `db:"TIPO_DOCUMENTO"`
	DocumentNumber           string     `db:"NUM_DOC_IDENTIDAD"`
	DocumentIssueDate        *time.Time `db:"DBU_FCH_EXD_DOC_IDENTIDAD"`
	Gender                   string     `db:"DBU_GENERO"`
	BirthDate                *time.Time `db:"DBU_FCH_NACIMIENTO"`
	BirthMunicipalityID      *int64      `db:"MPO_ID_NACIMIENTO"`
	BirthMunicipalityName    string     `db:"MPO_NOMBRE_NACIMIENTO"`
	MilitaryCard             int8       `db:"DBU_LIB_MILITAR"`
	MaritalStatus            string     `db:"DBU_ESTADO_CIVIL"`
	Stratum                  int8       `db:"DBU_ESTRATO"`
	BloodType                string     `db:"DBU_TIPO_SANGRE"`
	EpsAffiliated            string     `db:"DBU_AFILIADO_EPS"`
	EpsName                  string     `db:"DBU_EPS"`
	EmergencyContactName     string     `db:"DBU_NOMBRE_CONTACTO"`
	EmergencyPhoneLandline   string     `db:"DBU_TEL_FIJO_CONTACTO"`
	EmergencyPhoneMobile     string     `db:"DBU_TEL_MOVIL_CONTACTO"`
	EmergencyRelationship    string     `db:"DBU_PARENTESCO_CONTACTO"`
	EmergencyCompany         string     `db:"DBU_EMPRESA_CONTACTO"`
	EmergencyBirthDate       *time.Time `db:"DBU_FCH_NACIMIENTO_CONTACTO"`
	IcfesScore               float64    `db:"DBU_PUNTAJE_ICFES"`
	IsTechnicalHighSchool    string     `db:"DBU_ES_MEDIA_TECNICA"`
	DocExpMunicipalityID     *int64      `db:"MPO_ID_EXP_DOC_IDENTIDAD"`
	ShareInfoWithSena        string     `db:"DBU_COMPARTIR_INFO_SENA"`
	IsSisben                 string     `db:"DBU_SISBEN"`
	SisbenScore              int8       `db:"DBU_SISBEN_CALIFICA"`
	ChildrenCount            int8       `db:"DBU_HIJOS"`
	Ess                      string     `db:"DBU_ESS"`
	MilitaryCardNumber       string     `db:"DBU_NRO_LIB_MILITAR"`
	ExpirationDate           *time.Time `db:"DBU_FCH_VENCIMIENTO"`
  }
  

func (BasicDataUserModel) TableName() string {
	return "COMUN.DATOS_BASICOS_USUARIO"
}
