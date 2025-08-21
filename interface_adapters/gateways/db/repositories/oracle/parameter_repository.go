package oracle

import (
	"context"

	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/oracle"
	oracleDataSource "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/oracle"
)

// ParameterRepository maneja las operaciones de parámetros en Oracle
type ParameterRepository struct {
	parameterDataSource oracleDataSource.ParameterDataSource
}

// NewParameterRepository crea una nueva instancia de ParameterRepository
func NewParameterRepository(parameterDataSource oracleDataSource.ParameterDataSource) ParameterRepository {
	return ParameterRepository{
		parameterDataSource: parameterDataSource,
	}
}

// GetParameterByName obtiene un parámetro por su nombre
//
// :param ctx: contexto de la operación
// :param parameterName: nombre del parámetro a buscar
// :return: ParameterModel y error
func (p ParameterRepository) GetParameterByName(ctx context.Context, parameterName string) (*oracle.ParameterModel, error) {
	return p.parameterDataSource.GetParameterByName(ctx, parameterName)
}

// GetParameterValueByName obtiene solo el valor de un parámetro por su nombre
//
// :param ctx: contexto de la operación
// :param parameterName: nombre del parámetro a buscar
// :return: valor del parámetro como string y error
func (p ParameterRepository) GetParameterValueByName(ctx context.Context, parameterName string) (string, error) {
	return p.parameterDataSource.GetParameterValueByName(ctx, parameterName)
}
