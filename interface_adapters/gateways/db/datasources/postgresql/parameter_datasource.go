package postgresql

import (
	"context"
	"fmt"

	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/postgresql"
	"gorm.io/gorm"
)

// ParameterDataSource maneja las operaciones de la tabla PARAMETRO en PostgreSQL
type ParameterDataSource struct {
	connection *gorm.DB
}

// NewParameterDataSource crea una nueva instancia de ParameterDataSource
func NewParameterDataSource(connection *gorm.DB) ParameterDataSource {
	return ParameterDataSource{
		connection: connection,
	}
}

// GetParameterByName obtiene un parámetro por su nombre
//
// :param ctx: contexto de la operación
// :param parameterName: nombre del parámetro a buscar
// :return: ParameterModel y error
func (p ParameterDataSource) GetParameterByName(ctx context.Context, parameterName string) (*postgresql.ParameterModel, error) {
	var parameter postgresql.ParameterModel
	err := p.connection.WithContext(ctx).Where("\"PAR_NOMBRE\" = ?", parameterName).First(&parameter).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("parámetro '%s' no encontrado", parameterName)
		}
		return nil, fmt.Errorf("error consultando parámetro '%s': %w", parameterName, err)
	}

	return &parameter, nil
}

// GetParameterValueByName obtiene solo el valor de un parámetro por su nombre
//
// :param ctx: contexto de la operación
// :param parameterName: nombre del parámetro a buscar
// :return: valor del parámetro como string y error
func (p ParameterDataSource) GetParameterValueByName(ctx context.Context, parameterName string) (string, error) {
	var value string
	err := p.connection.WithContext(ctx).Model(&postgresql.ParameterModel{}).
		Where("par_nombre = ?", parameterName).
		Pluck("par_valor", &value).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("parámetro '%s' no encontrado", parameterName)
		}
		return "", fmt.Errorf("error consultando valor del parámetro '%s': %w", parameterName, err)
	}

	if len(value) == 0 {
		return "", fmt.Errorf("parámetro '%s' no encontrado", parameterName)
	}

	return value, nil
}
