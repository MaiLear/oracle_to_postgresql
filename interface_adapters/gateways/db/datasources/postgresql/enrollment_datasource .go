package postgresql

import (
	"context"
	"errors"
	cockroachdbErrors "github.com/cockroachdb/errors"
	errorDbConnector "gitlab.com/sofia-plus/go_db_connectors/errors"
	postgresqlPort "gitlab.com/sofia-plus/go_db_connectors/ports/repositories/postgresql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/constants"
	postgresqlModels "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/postgresql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/postgresql/dto"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
	"gorm.io/gorm"
)

type EnrollmentDataSource struct {
	connection *gorm.DB
	db         postgresqlPort.PostgresqlRepository
}

func NewEnrollmentDataSource(db postgresqlPort.PostgresqlRepository, connection *gorm.DB) EnrollmentDataSource {
	return EnrollmentDataSource{
		connection: connection,
		db:         db,
	}
}

func (e EnrollmentDataSource) GetEventVenue(ctx context.Context, charaSheetId int) (vanueEvent []*postgresqlModels.VanueEventModel, err error) {
	where := e.db.Where(`"FIC_ID"`, "=", charaSheetId)
	if err = e.db.SelectAll(&vanueEvent, where); err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInPostgres) {
			return vanueEvent, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound, "infra: registros de evento sede no encontrados en postgresql:  %s", err.Error())
		}
		return vanueEvent, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo los eventos sede")
	}
	return vanueEvent, nil
}

func (e EnrollmentDataSource) GetPendingRecords(ctx context.Context) (enrollWithProgramAndCourse []dto.EnrollmentWithProgramAndCourse, err error) {
	result := e.connection.Table(`inscription."INGRESO_ASPIRANTE" AS IA`).
		Joins(`INNER JOIN program."FICHA_CARACTERIZACION" AS FC ON FC."FIC_ID" = IA."FIC_ID"`).
		Joins(`INNER JOIN program."PROGRAMA_FORMACION" AS PRF ON PRF."PRF_ID" = FC."PRF_ID"`).
		Where(`IA."state" = ?`, "pending").
		Select(`
		    IA."ING_ID",
            IA."ING_ID_PG",
            IA."ING_PERIODO",
            IA."NIS",
            IA."FIC_ID",
            IA."ING_ESTADO",
            IA."ING_FCH_REGISTRO",
            IA."NIS_FUN_REGISTRO",
            IA."ING_APLICO_CONVENIO",
            IA."ING_NUMERO_CONVENIO",
            IA."ING_DECRETO_REINTEGRADO",
            IA."ING_OBS_REINTEGRADO",
            IA."ING_PORCENTAJE_PONDERACION",
            IA."ING_PUNTAJE_TOTAL",
            IA."ING_PRIORIDAD_VIRTUAL",
            IA."FTE_ID_ACCESO_PREF",
            IA."state",
			IA."REQUEST_DATA",
			IA."number_attemps",
            FC."SOS_ID",
            PRF."PRF_TIPO_PROGRAMA",
			PRF."PRF_DENOMINACION"
		`).
		Scan(&enrollWithProgramAndCourse)	

	if result.Error != nil {
		return nil, cockroachdbErrors.Wrap(result.Error, "infra: ocurrió un error obteniendo las inscripciones con su ficha y programa de formacion")
	}

	if result.RowsAffected == 0 {
		return nil, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound, "infra: no se encontraron inscripciones con sus relaciones")
	}
	return enrollWithProgramAndCourse, nil
}

// func (e EnrollmentDataSource) GetPendingRecords(ctx context.Context) (enrollmentWithPendingState []*postgresqlModels.EnrollmentModel, err error) {
// 	where := e.db.Where(`"state"`, "=", "pending")
// 	if err = e.db.SelectAll(&enrollmentWithPendingState, where); err != nil {
// 		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInPostgres) {
// 			return enrollmentWithPendingState, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound, "infra: registros de ingreso aspirante pendientes no encontrados en postgresql:  %s", err.Error())
// 		}
// 		return enrollmentWithPendingState, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo los registros pendientes de ingreso aspirante")
// 	}
// 	return enrollmentWithPendingState, nil
// }

func (e EnrollmentDataSource) GetSynchronizedRecords(ctx context.Context) (enrollmentWithSyncedState []*postgresqlModels.EnrollmentModel, err error) {
	firstWhere := e.db.Where(`"state"`, "=", "synced")
	secondWhere := e.db.Where(`"ING_ESTADO"`, "=", 1)
	if err = e.db.SelectAll(&enrollmentWithSyncedState, firstWhere, secondWhere); err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInPostgres) {
			return enrollmentWithSyncedState, cockroachdbErrors.Wrapf(internalErrors.ErrNotFound, "infra: registros de ingreso aspirante sincronizados no encontrados en postgresql:  %s", err.Error())
		}
		return enrollmentWithSyncedState, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo los registros pendientes de ingreso aspirante")
	}
	return enrollmentWithSyncedState, nil
}

// TODO: falta la columna ING_ID_PG en la tabla de postgres para hacer este metodo
func (u EnrollmentDataSource) UpdateSyncStatusWithID(idFromOra, idFromPos int) error {
	result := u.connection.Exec(`UPDATE inscription."INGRESO_ASPIRANTE" SET "state" = $1,"ING_ID" = $2 WHERE "ING_ID_PG"= $3`, constants.SyncedStatus, idFromOra, idFromPos)
	if result.Error != nil {
		return cockroachdbErrors.Wrap(result.Error, "infra: ocurrio un error modificando el estado y clave original del  ingreso aspirante")
	}
	return nil
}

func (e EnrollmentDataSource) UpdateEnrollment(ctx context.Context, enrollmentModel postgresqlModels.EnrollmentModel) error {
	result := e.connection.Model(&enrollmentModel).Updates(enrollmentModel)
	if result.Error != nil {
		return cockroachdbErrors.Wrapf(result.Error, "infra: ocurrio un error actualizacion el ingreso aspirante en postgresql: %s", result.Error.Error())
	}

	return nil
}

func (e EnrollmentDataSource) UpdateEnrollmentStatusError(ctx context.Context, ingIDPG int) error {
	result := e.connection.
		Model(&postgresqlModels.EnrollmentModel{}).
		Where(`"ING_ID_PG" = ?`, ingIDPG).
		Update("state", "error")

	if result.Error != nil {
		return cockroachdbErrors.Wrapf(
			result.Error,
			"infra: ocurrió un error actualizando el status a 'error' para el ingreso con ING_ID_PG=%d: %s",
			ingIDPG, result.Error.Error(),
		)
	}
	if result.RowsAffected == 0 {
		return cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,
			"infra: no se encontró ningún registro con ING_ID_PG=%d para actualizar a 'error'",
			ingIDPG,
		)
	}

	return nil
}

func (e EnrollmentDataSource) SetNumberAttemps(ctx context.Context, ingIDPG int) (err error) {
	result := e.connection.
		Model(&postgresqlModels.EnrollmentModel{}).
		Where(`"ING_ID_PG" = ?`, ingIDPG).
		Update("number_attemps",gorm.Expr("number_attemps + ?", 1))

	if result.Error != nil {
		return cockroachdbErrors.Wrapf(
			result.Error,
			"infra: ocurrió un error actualizando el numero de intentos con ING_ID_PG=%d: %s",
			ingIDPG, result.Error.Error(),
		)
	}
	if result.RowsAffected == 0 {
		return cockroachdbErrors.Wrapf(internalErrors.ErrNotFound,
			"infra: no se encontró ningún registro con ING_ID_PG=%d para actualizar el numero de intentos",
			ingIDPG,
		)
	}

	return nil
}

