package postgresql

import (
	"context"
	"fmt"
	"time"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/postgresql"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
	"gorm.io/gorm"
)

type DocumentTypeDatasource struct {
	db *gorm.DB
}

func NewDocumentTypeDatasource(db *gorm.DB) DocumentTypeDatasource {
	return DocumentTypeDatasource{db: db}
}

func (d DocumentTypeDatasource) GetValidTypesForRegistration(ctx context.Context) ([]entities.DocumentType, error) {
	var models []postgresql.DocumentTypeModel
	err := d.db.WithContext(ctx).Where("tdoc_registro = ?", "1").Order("tdoc_orden").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("error getting valid document types for registration: %w", err)
	}

	var entities []entities.DocumentType
	for _, model := range models {
		entities = append(entities, model.ToEntity())
	}

	return entities, nil
}

func (d DocumentTypeDatasource) GetByType(ctx context.Context, documentType string) (*entities.DocumentType, error) {
	var model postgresql.DocumentTypeModel
	err := d.db.WithContext(ctx).Where("tdoc_tipo = ?", documentType).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound, "infra: no se encontro el tipo de documento")
		}
		return nil, cockroachdbErrors.Wrap(err, "infra: ocurrio un error encontrando el tipo de documento")
	}

	entity := model.ToEntity()
	return &entity, nil
}

// GetPPTValidityYears obtiene los años de vigencia del PPT desde la tabla TIPO_DOCUMENTO
//
// :param ctx: contexto de la operación
// :return: años de vigencia del PPT y error
func (d DocumentTypeDatasource) GetPPTValidityYears(ctx context.Context) (int, error) {
	var duration int
	err := d.db.WithContext(ctx).Model(&postgresql.DocumentTypeModel{}).
		Where("tdoc_tipo = ?", "PPT").
		Pluck("tdoc_duracion", &duration).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no se encontró la configuración de duración para el documento PPT")
		}
		return 0, fmt.Errorf("error consultando duración del PPT: %w", err)
	}

	if duration == 0 {
		return 0, fmt.Errorf("no se encontró la configuración de duración para el documento PPT")
	}

	return duration, nil
}

// GetPPTExpirationDate obtiene la fecha de vencimiento del PPT desde la tabla TIPO_DOCUMENTO
//
// :param ctx: contexto de la operación
// :return: fecha de vencimiento del PPT y error
func (d DocumentTypeDatasource) GetPPTExpirationDate(ctx context.Context) (*time.Time, error) {
	var model postgresql.DocumentTypeModel
	err := d.db.WithContext(ctx).Where("tdoc_tipo = ?", "PPT").First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no se encontró la configuración de fecha de vencimiento para el documento PPT")
		}
		return nil, fmt.Errorf("error consultando fecha de vencimiento del PPT: %w", err)
	}

	if !model.ExpirationDate.Valid {
		return nil, fmt.Errorf("no se encontró la configuración de fecha de vencimiento para el documento PPT")
	}

	return &model.ExpirationDate.Time, nil
}

func (d DocumentTypeDatasource) GetAll(ctx context.Context) ([]entities.DocumentType, error) {
	var models []postgresql.DocumentTypeModel
	err := d.db.WithContext(ctx).Order("tdoc_orden").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("error getting all document types: %w", err)
	}

	var entities []entities.DocumentType
	for _, model := range models {
		entities = append(entities, model.ToEntity())
	}

	return entities, nil
}
