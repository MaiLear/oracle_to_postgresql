package oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	cockroachdbErrors "github.com/cockroachdb/errors"
	errorDbConnector "gitlab.com/sofia-plus/go_db_connectors/errors"
	oraclePort "gitlab.com/sofia-plus/go_db_connectors/ports/repositories/oracle"
	oracleModels "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/oracle"
	//"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/oracle/dto"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
	//"gitlab.com/sofia-plus/go_db_connectors/repositories/oracle"
)

type EnrollmentDataSource struct {
	connection *sql.DB
	db         oraclePort.OracleRepository
}

func NewEnrollmentDataSource(db oraclePort.OracleRepository, connection *sql.DB) EnrollmentDataSource {
	return EnrollmentDataSource{
		connection: connection,
		db:         db,
	}
}

// func (e EnrollmentDataSource) SaveEnrollment(ctx context.Context, enrollmentDto dto.EnrollmentDto) (id,nis int, err error) {
// 	var returned  = struct{
// 		ID int `db:"ING_ID"`
// 		Nis int `db:"NIS"`
// 	}{}
// 	err = e.db.Insert(enrollmentDto, &returned)
// 	if err != nil{
// 		fmt.Printf("\n OCURRIO UN ERROR AL INSERTAR EN ORCLE LA ENTIDAD ES\n")
// 		fmt.Printf("\n %+v \n",enrollmentDto)
// 		err = cockroachdbErrors.Wrap(err,"infra: ocurrio un error insertando el registro de inscripcion aspirante en oracle")
// 		return
// 	}
// 	id = returned.ID
// 	nis = returned.Nis
// 	fmt.Printf("\nING_ID oracle %d y NIS ORACLE %d\n",returned.ID,returned.Nis)
// 	return
// }


func (e EnrollmentDataSource) GetEnrollmentByNis(ctx context.Context, nis int) (enrollment *oracleModels.EnrollmentModel, err error) {
	model := new(oracleModels.EnrollmentModel)
	where := e.db.Where("NIS = :1", nis)
	result, err := e.db.Select(model, where)
	if err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInOracle) {
			return enrollment, cockroachdbErrors.Wrap(internalErrors.ErrNotFound, fmt.Sprintf("infra: no se encontro la inscripcion por el nis %s", err.Error()))
		}
		return enrollment, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo el curso en datasource")
	}
	enrollment, ok := result.(*oracleModels.EnrollmentModel)
	if !ok {
		return nil, cockroachdbErrors.New("infra: error de conversión de tipo al obtener las inscripciones aspirante")
	}
	return enrollment, nil
}


//TODO: MIRAR QUE FUNCINE BIEN
func (e EnrollmentDataSource) GetAllEnrollmentsByIds(ctx context.Context, ids []int) (enrollments []*oracleModels.EnrollmentModel, err error) {
	if len(ids) == 0 {
		err = cockroachdbErrors.New("infra: no hay ids para hacer la busquedad de la inscripcion aspirante en oracle")
		return 
	}
	placeholders := make([]string, len(ids))
	params := make([]any, len(ids))
	for i, val := range ids {
		placeholders[i] = fmt.Sprintf(":%d", i+1) // Oracle usa :1, :2, ...
		params[i] = val
	}
	model := new(oracleModels.EnrollmentModel)
	whereClause := fmt.Sprintf("ING_ID IN (%s)", strings.Join(placeholders, ", "))
	where := e.db.Where(whereClause,params...)
	result, err := e.db.SelectAll(model, where)
	if err != nil {
		if errors.Is(err, errorDbConnector.ErrRecordNotFoundInOracle) {
			return enrollments, cockroachdbErrors.Wrap(internalErrors.ErrNotFound, fmt.Sprintf("infra: no se encontro el ingreso aspirante por el id en oracle %s",err.Error()))
		}
		return enrollments, cockroachdbErrors.Wrap(err, "infra: ocurrio un error obteniendo el ingreso aspirante")
	}
	enrollments, ok := result.([]*oracleModels.EnrollmentModel)
	if !ok {
		return nil, cockroachdbErrors.New("infra: error de conversión de tipo al obtener las inscripciones aspirante")
	}
	return enrollments, nil
}



