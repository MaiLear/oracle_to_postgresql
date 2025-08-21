package oracle

import (
	"context"
	"fmt"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	oracleDataSource "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/oracle"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
)

type EnrollmentRepository struct {
	datasourceEnrollmentStored oracleDataSource.EnrrollmentStoredProceduteDatasource
	datasource oracleDataSource.EnrollmentDataSource
}

func NewEnrollmentRepository(datasource oracleDataSource.EnrollmentDataSource,datasourceEnrollmentStored oracleDataSource.EnrrollmentStoredProceduteDatasource) EnrollmentRepository {
	return EnrollmentRepository{
		datasource: datasource,
		datasourceEnrollmentStored: datasourceEnrollmentStored,
	}
}

func (u EnrollmentRepository) SaveEnrollment(ctx context.Context, enrollmentDomain entities.Enrollment) (messageErr string,err error) {
	//enrollmentDto := mappers.FromEnrollmentDomainToDto(enrollmentDomain)
	messageErr,err = u.datasourceEnrollmentStored.Enroll(ctx, enrollmentDomain)
	if err != nil {
		return messageErr,cockroachdbErrors.WithStack(err)
	}
	return "",nil
}

func (u EnrollmentRepository) SaveComplementaryEnrollment(ctx context.Context, enrollmentDomain entities.Enrollment) (messageErr string,err error) {
	//enrollmentDto := mappers.FromEnrollmentDomainToDto(enrollmentDomain)
	messageErr,err = u.datasourceEnrollmentStored.EnrollComplementary(ctx, enrollmentDomain)
	if err != nil {
		fmt.Printf("MENSAJE DE ERROR EN LA LINEA 38 EN REPOSITORIO")
		return messageErr,cockroachdbErrors.WithStack(err)
	}
	return "",nil
}

func (e EnrollmentRepository) GetAllEnrollmentsByIds(ctx context.Context, ids []int) (enrollmentDomain []entities.Enrollment, err error) {
	enrollmentsModels, err := e.datasource.GetAllEnrollmentsByIds(ctx, ids)
	if err != nil {
		return enrollmentDomain, cockroachdbErrors.WithStack(err)
	}
	enrollmentDomain = mappers.FromEnrollModelsOraToDomains(enrollmentsModels)
	return enrollmentDomain, nil
}


func (e EnrollmentRepository) GetEnrollmentByNis(ctx context.Context, nis int) (enrollmentDomain entities.Enrollment, err error) {
	enrollmentModel, err := e.datasource.GetEnrollmentByNis(ctx, nis)
	if err != nil {
		return enrollmentDomain, cockroachdbErrors.WithStack(err)
	}
	enrollmentDomain = mappers.FromOraEnrrollModelToDomain(enrollmentModel)
	return enrollmentDomain, nil
}
