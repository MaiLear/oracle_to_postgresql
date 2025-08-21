package oracle

import (
	"database/sql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	"time"
)

type DocumentTypeModel struct {
	Type                    string       `db:"TDOC_TIPO"`
	Name                    string       `db:"TDOC_NOMBRE"`
	Registration            string       `db:"TDOC_REGISTRO"`
	Inscription             string       `db:"TDOC_INSCRIPCION"`
	Login                   string       `db:"TDOC_LOGIN"`
	Duration                int          `db:"TDOC_DURACION"`
	ExpirationDate          sql.NullTime `db:"TDOC_FCH_VENCIMIENTO"`
	Order                   int          `db:"TDOC_ORDEN"`
	Consultations           string       `db:"TDOC_CONSULTAS"`
	TypeIdentificationCall  string       `db:"TDOC_TIPO_IDENTIF_CALL"`
	TypeIdentificationPers  string       `db:"TDOC_TIPO_IDENTIF_PERS"`
	TypeIdentification      string       `db:"TDOC_TIPO_IDENTIF"`
	TypeCompany             string       `db:"TDOC_TIPO_EMPRESA"`
	InsApoyosSos            string       `db:"TDOC_INS_APOYOS_SOS"`
	ApoyosFic               string       `db:"TDOC_APOYOS_FIC"`
	TypeIdentificationDB    string       `db:"TDOC_TIPO_IDENTIF_DAT_BAS"`
	CreateCompany           string       `db:"TDOC_CREAR_EMPRESA"`
	TypeIdentificationPopup string       `db:"TDOC_TIPO_IDENTIF_POPUP"`
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
	return "COMUN.TIPO_DOCUMENTO"
}
