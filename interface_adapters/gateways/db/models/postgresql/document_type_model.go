package postgresql

import (
	"database/sql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	"time"
)

type DocumentTypeModel struct {
	Type                    string       `db:"tdoc_tipo" json:"type"`
	Name                    string       `db:"tdoc_nombre" json:"name"`
	Registration            string       `db:"tdoc_registro" json:"registration"`
	Inscription             string       `db:"tdoc_inscripcion" json:"inscription"`
	Login                   string       `db:"tdoc_login" json:"login"`
	Duration                int          `db:"tdoc_duracion" json:"duration"`
	ExpirationDate          sql.NullTime `db:"tdoc_fch_vencimiento" json:"expiration_date"`
	Order                   int          `db:"tdoc_orden" json:"order"`
	Consultations           string       `db:"tdoc_consultas" json:"consultations"`
	TypeIdentificationCall  string       `db:"tdoc_tipo_identif_call" json:"type_identification_call"`
	TypeIdentificationPers  string       `db:"tdoc_tipo_identif_pers" json:"type_identification_pers"`
	TypeIdentification      string       `db:"tdoc_tipo_identif" json:"type_identification"`
	TypeCompany             string       `db:"tdoc_tipo_empresa" json:"type_company"`
	InsApoyosSos            string       `db:"tdoc_ins_apoyos_sos" json:"ins_apoyos_sos"`
	ApoyosFic               string       `db:"tdoc_apoyos_fic" json:"apoyos_fic"`
	TypeIdentificationDB    string       `db:"tdoc_tipo_identif_dat_bas" json:"type_identification_db"`
	CreateCompany           string       `db:"tdoc_crear_empresa" json:"create_company"`
	TypeIdentificationPopup string       `db:"tdoc_tipo_identif_popup" json:"type_identification_popup"`
}

func (m DocumentTypeModel) ToEntity() entities.DocumentType {
	var expirationDate *time.Time
	if m.ExpirationDate.Valid {
		expirationDate = &m.ExpirationDate.Time
	}

	return entities.DocumentType{
		Type:                    m.Type,
		Name:                    m.Name,
		AllowsRegistration:      m.Registration == "1",
		AllowsInscription:       m.Inscription == "1",
		AllowsLogin:             m.Login == "1",
		Duration:                m.Duration,
		ExpirationDate:          expirationDate,
		Order:                   m.Order,
		AllowsConsultations:     m.Consultations == "1",
		TypeIdentificationCall:  m.TypeIdentificationCall == "1",
		TypeIdentificationPers:  m.TypeIdentificationPers == "1",
		TypeIdentification:      m.TypeIdentification == "1",
		TypeCompany:             m.TypeCompany == "1",
		InsApoyosSos:            m.InsApoyosSos == "1",
		ApoyosFic:               m.ApoyosFic == "1",
		TypeIdentificationDB:    m.TypeIdentificationDB == "1",
		CreateCompany:           m.CreateCompany == "1",
		TypeIdentificationPopup: m.TypeIdentificationPopup == "1",
	}
}

func (m DocumentTypeModel) TableName() string {
	return "common.\"TIPO_DOCUMENTO\""
}
