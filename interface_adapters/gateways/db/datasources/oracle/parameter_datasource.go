package oracle

import (
	"context"
	"database/sql"
	"fmt"

	oracleModels "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/oracle"
)

// ParameterDataSource maneja las operaciones de la tabla COMUN.PARAMETRO en Oracle
type ParameterDataSource struct {
	connection *sql.DB
}

// NewParameterDataSource crea una nueva instancia de ParameterDataSource
func NewParameterDataSource(connection *sql.DB) ParameterDataSource {
	return ParameterDataSource{
		connection: connection,
	}
}

// GetParameterByName obtiene un parámetro por su nombre
//
// :param ctx: contexto de la operación
// :param parameterName: nombre del parámetro a buscar
// :return: ParameterModel y error
func (p ParameterDataSource) GetParameterByName(ctx context.Context, parameterName string) (*oracleModels.ParameterModel, error) {
	query := `
		SELECT 
			PAR_ID,
			PAR_NOMBRE,
			PAR_TIPO,
			PAR_DESCRIPCION,
			PAR_VALOR,
			PAQ_ID,	
			PAR_FCH_REGISTRO,
			PAR_FCH_MODIF
		FROM COMUN.PARAMETRO 
		WHERE PAR_NOMBRE = :parameterName
	`

	var parameter oracleModels.ParameterModel
	err := p.connection.QueryRowContext(ctx, query, sql.Named("parameterName", parameterName)).Scan(
		&parameter.ParId,
		&parameter.ParNombre,
		&parameter.ParTipo,
		&parameter.ParDescripcion,
		&parameter.ParValor,
		&parameter.PaqId,
		&parameter.ParFchRegistro,
		&parameter.ParFchModif,
	)

	if err != nil {
		if err == sql.ErrNoRows {
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
	query := `
		SELECT PAR_VALOR
		FROM COMUN.PARAMETRO 
		WHERE PAR_NOMBRE = :parameterName
	`

	var value string
	err := p.connection.QueryRowContext(ctx, query, sql.Named("parameterName", parameterName)).Scan(&value)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("parámetro '%s' no encontrado", parameterName)
		}
		return "", fmt.Errorf("error consultando valor del parámetro '%s': %w", parameterName, err)
	}

	return value, nil
}
