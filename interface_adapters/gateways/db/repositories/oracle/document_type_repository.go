package oracle

import (
	"context"
	"time"

	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	oraclePort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/oracle"
)

type DocumentTypeRepository struct {
	datasource oraclePort.DocumentTypeDatasource
}

func NewDocumentTypeRepository(datasource oraclePort.DocumentTypeDatasource) DocumentTypeRepository {
	return DocumentTypeRepository{datasource: datasource}
}

func (r DocumentTypeRepository) GetValidTypesForRegistration(ctx context.Context) ([]entities.DocumentType, error) {
	return r.datasource.GetValidTypesForRegistration(ctx)
}

func (r DocumentTypeRepository) GetByType(ctx context.Context, documentType string) (*entities.DocumentType, error) {
	return r.datasource.GetByType(ctx, documentType)
}

func (r DocumentTypeRepository) GetAll(ctx context.Context) ([]entities.DocumentType, error) {
	return r.datasource.GetAll(ctx)
}

// GetPPTValidityYears obtiene los años de vigencia del PPT
//
// :param ctx: contexto de la operación
// :return: años de vigencia del PPT y error
func (r DocumentTypeRepository) GetPPTValidityYears(ctx context.Context) (int, error) {
	return r.datasource.GetPPTValidityYears(ctx)
}

// GetPPTExpirationDate obtiene la fecha de vencimiento del PPT desde Oracle
//
// :param ctx: contexto de la operación
// :return: fecha de vencimiento del PPT y error
func (r DocumentTypeRepository) GetPPTExpirationDate(ctx context.Context) (*time.Time, error) {
	return r.datasource.GetPPTExpirationDate(ctx)
}
