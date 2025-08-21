package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/oracle"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
)

type DocumentTypeDatasource struct {
	db *sql.DB
}

func NewDocumentTypeDatasource(db *sql.DB) DocumentTypeDatasource {
	return DocumentTypeDatasource{db: db}
}

func (d DocumentTypeDatasource) GetValidTypesForRegistration(ctx context.Context) ([]entities.DocumentType, error) {
	query := `
		SELECT TDOC_TIPO, TDOC_NOMBRE, TDOC_REGISTRO, TDOC_INSCRIPCION, TDOC_LOGIN, 
			   TDOC_DURACION, TDOC_FCH_VENCIMIENTO, TDOC_ORDEN, TDOC_CONSULTAS,
			   TDOC_TIPO_IDENTIF_CALL, TDOC_TIPO_IDENTIF_PERS, TDOC_TIPO_IDENTIF,
			   TDOC_TIPO_EMPRESA, TDOC_INS_APOYOS_SOS, TDOC_APOYOS_FIC,
			   TDOC_TIPO_IDENTIF_DAT_BAS, TDOC_CREAR_EMPRESA, TDOC_TIPO_IDENTIF_POPUP
		FROM COMUN.TIPO_DOCUMENTO 
		WHERE TDOC_REGISTRO = '1'
		ORDER BY TDOC_ORDEN`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting valid document types for registration: %w", err)
	}
	defer rows.Close()

	var entities []entities.DocumentType
	for rows.Next() {
		var model oracle.DocumentTypeModel
		err := rows.Scan(
			&model.Type, &model.Name, &model.Registration, &model.Inscription, &model.Login,
			&model.Duration, &model.ExpirationDate, &model.Order, &model.Consultations,
			&model.TypeIdentificationCall, &model.TypeIdentificationPers, &model.TypeIdentification,
			&model.TypeCompany, &model.InsApoyosSos, &model.ApoyosFic,
			&model.TypeIdentificationDB, &model.CreateCompany, &model.TypeIdentificationPopup,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning document type row: %w", err)
		}
		entities = append(entities, model.ToEntity())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating document type rows: %w", err)
	}

	return entities, nil
}

func (d DocumentTypeDatasource) GetByType(ctx context.Context, documentType string) (*entities.DocumentType, error) {
	query := `
		SELECT TDOC_TIPO, TDOC_NOMBRE, TDOC_REGISTRO, TDOC_INSCRIPCION, TDOC_LOGIN, 
			   TDOC_DURACION, TDOC_FCH_VENCIMIENTO, TDOC_ORDEN, TDOC_CONSULTAS,
			   TDOC_TIPO_IDENTIF_CALL, TDOC_TIPO_IDENTIF_PERS, TDOC_TIPO_IDENTIF,
			   TDOC_TIPO_EMPRESA, TDOC_INS_APOYOS_SOS, TDOC_APOYOS_FIC,
			   TDOC_TIPO_IDENTIF_DAT_BAS, TDOC_CREAR_EMPRESA, TDOC_TIPO_IDENTIF_POPUP
		FROM COMUN.TIPO_DOCUMENTO 
		WHERE TDOC_TIPO = :1`

	var model oracle.DocumentTypeModel
	row := d.db.QueryRowContext(ctx, query, documentType)
	err := row.Scan(
		&model.Type, &model.Name, &model.Registration, &model.Inscription, &model.Login,
		&model.Duration, &model.ExpirationDate, &model.Order, &model.Consultations,
		&model.TypeIdentificationCall, &model.TypeIdentificationPers, &model.TypeIdentification,
		&model.TypeCompany, &model.InsApoyosSos, &model.ApoyosFic,
		&model.TypeIdentificationDB, &model.CreateCompany, &model.TypeIdentificationPopup,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound, "infra: no se encontro el tipo de documento")
		}
		return nil, cockroachdbErrors.Wrap(err, "infra: ocurrio un error encontrando el tipo de documento")
	}

	entity := model.ToEntity()
	return &entity, nil
}

func (d DocumentTypeDatasource) GetAll(ctx context.Context) ([]entities.DocumentType, error) {
	query := `
		SELECT TDOC_TIPO, TDOC_NOMBRE, TDOC_REGISTRO, TDOC_INSCRIPCION, TDOC_LOGIN, 
			   TDOC_DURACION, TDOC_FCH_VENCIMIENTO, TDOC_ORDEN, TDOC_CONSULTAS,
			   TDOC_TIPO_IDENTIF_CALL, TDOC_TIPO_IDENTIF_PERS, TDOC_TIPO_IDENTIF,
			   TDOC_TIPO_EMPRESA, TDOC_INS_APOYOS_SOS, TDOC_APOYOS_FIC,
			   TDOC_TIPO_IDENTIF_DAT_BAS, TDOC_CREAR_EMPRESA, TDOC_TIPO_IDENTIF_POPUP
		FROM COMUN.TIPO_DOCUMENTO 
		ORDER BY TDOC_ORDEN`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting all document types: %w", err)
	}
	defer rows.Close()

	var entities []entities.DocumentType
	for rows.Next() {
		var model oracle.DocumentTypeModel
		err := rows.Scan(
			&model.Type, &model.Name, &model.Registration, &model.Inscription, &model.Login,
			&model.Duration, &model.ExpirationDate, &model.Order, &model.Consultations,
			&model.TypeIdentificationCall, &model.TypeIdentificationPers, &model.TypeIdentification,
			&model.TypeCompany, &model.InsApoyosSos, &model.ApoyosFic,
			&model.TypeIdentificationDB, &model.CreateCompany, &model.TypeIdentificationPopup,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning document type row: %w", err)
		}
		entities = append(entities, model.ToEntity())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating document type rows: %w", err)
	}

	return entities, nil
}

// GetPPTValidityYears obtiene los años de vigencia del PPT desde la tabla TIPO_DOCUMENTO
//
// :param ctx: contexto de la operación
// :return: años de vigencia del PPT y error
func (d DocumentTypeDatasource) GetPPTValidityYears(ctx context.Context) (int, error) {
	query := `
		SELECT TDOC_DURACION
		FROM COMUN.TIPO_DOCUMENTO 
		WHERE TDOC_TIPO = 'PPT'
	`

	var duration int
	err := d.db.QueryRowContext(ctx, query).Scan(&duration)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no se encontró la configuración de duración para el documento PPT")
		}
		return 0, fmt.Errorf("error consultando duración del PPT: %w", err)
	}

	return duration, nil
}

// GetPPTExpirationDate obtiene la fecha de vencimiento del PPT desde la tabla TIPO_DOCUMENTO
//
// :param ctx: contexto de la operación
// :return: fecha de vencimiento del PPT y error
func (d DocumentTypeDatasource) GetPPTExpirationDate(ctx context.Context) (*time.Time, error) {
	query := `
		SELECT TDOC_FCH_VENCIMIENTO
		FROM COMUN.TIPO_DOCUMENTO 
		WHERE TDOC_TIPO = 'PPT'
	`

	var expirationDate sql.NullTime
	err := d.db.QueryRowContext(ctx, query).Scan(&expirationDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no se encontró la configuración de fecha de vencimiento para el documento PPT")
		}
		return nil, fmt.Errorf("error consultando fecha de vencimiento del PPT: %w", err)
	}

	if !expirationDate.Valid {
		return nil, fmt.Errorf("no se encontró la configuración de fecha de vencimiento para el documento PPT")
	}

	return &expirationDate.Time, nil
}
