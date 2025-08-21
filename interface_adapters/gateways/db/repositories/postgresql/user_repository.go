package postgresql

import (
	"context"
	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	dataSourcePort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/postgresql"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
)

// UserRepository implementa la interfaz de acceso a datos para Programas.
type UserRepository struct {
	datasource dataSourcePort.UserDataSource
}

func NewUserRepository(datasource dataSourcePort.UserDataSource) UserRepository {
	return UserRepository{datasource: datasource}
}

func (u UserRepository) GetPendingRecords(ctx context.Context) (domainUsers []entities.User, err error) {
	usersWithPendingState, err := u.datasource.GetPendingRecords(ctx)
	if err != nil {
		return domainUsers, cockroachdbErrors.WithStack(err)
	}
	domainUsers = mappers.FromUsersModelPosToDomain(usersWithPendingState)
	return domainUsers, nil
}

func (u UserRepository) UpdateEnrollmentStatusError(ctx context.Context, nisFromPos int) error {
	if err := u.datasource.UpdateEnrollmentStatusError(ctx, nisFromPos); err != nil {
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}

func (u UserRepository) SetNumberAttemps(ctx context.Context, nisFromPos int) (err error){
	if err := u.datasource.SetNumberAttemps(ctx,nisFromPos); err != nil{
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}


func (u UserRepository) UpdateSyncStatusWithID(nisFromOra,nisFromPost int) error {
	if err := u.datasource.UpdateSyncStatusWithID(nisFromOra,nisFromPost); err != nil {
		return cockroachdbErrors.WithStack(err)
	}
	return nil
}
