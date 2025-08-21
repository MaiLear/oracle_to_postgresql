package postgresql

import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	dataSourcePort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/postgresql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
)

// EnrollmentRepository implementa la interfaz de acceso a datos para Programas.
type EnrollmentRepository struct {
	datasource dataSourcePort.EnrollmentDataSource
}

func NewEnrollmentRepository(datasource dataSourcePort.EnrollmentDataSource) EnrollmentRepository {
	return EnrollmentRepository{datasource: datasource}
}

func (e EnrollmentRepository) GetPendingRecords(ctx context.Context) (domainEnrollment []entities.Enrollment, err error) {
	enrollmentWithPendingState, err := e.datasource.GetPendingRecords(ctx)
	if err != nil {
		return domainEnrollment, cockroachdbErrors.WithStack(err)
	}
	domainEnrollment = mappers.FromEnrollmentDtoDatasourceToDomain(enrollmentWithPendingState)
	return domainEnrollment, nil
}

func (e EnrollmentRepository) GetSynchronizedRecords(ctx context.Context) (domainEnrollment []entities.Enrollment, err error) {
	enrollmentWithSyncronizeState, err := e.datasource.GetSynchronizedRecords(ctx)
	if err != nil {
		return domainEnrollment, cockroachdbErrors.WithStack(err)
	}
	domainEnrollment = mappers.FromEnrollmentModelPosToDomain(enrollmentWithSyncronizeState)
	return domainEnrollment, nil
}

func (u EnrollmentRepository) UpdateSyncStatusWithID(idFromOra,idFroPos int) error {
	if err := u.datasource.UpdateSyncStatusWithID(idFromOra,idFroPos); err != nil {
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}

func (e EnrollmentRepository) UpdateEnrollment(ctx context.Context,enrollmentDomain entities.Enrollment)error{
	enrollmentModel := mappers.FromEnrollmentDomainToPosModel(enrollmentDomain)
	if err := e.datasource.UpdateEnrollment(ctx,enrollmentModel); err != nil{
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}

func (e EnrollmentRepository)GetEventVenue(ctx context.Context, charaSheetId int)(err error){
	_,err = e.datasource.GetEventVenue(ctx,charaSheetId)
	if err != nil{
		err = cockroachdbErrors.WithStack(err)
		return
	}
	return
}

func (e EnrollmentRepository) UpdateEnrollmentStatusError(ctx context.Context, ingIDPG int) error {
	if err := e.datasource.UpdateEnrollmentStatusError(ctx, ingIDPG); err != nil {
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}

func (e EnrollmentRepository) SetNumberAttemps(ctx context.Context, ingIDPG int) (err error){
	if err := e.datasource.SetNumberAttemps(ctx,ingIDPG); err != nil{
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}

