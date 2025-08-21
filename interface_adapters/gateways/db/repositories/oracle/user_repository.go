package oracle

import (
	"context"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/entities"
	oracleDataSource "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/datasources/oracle"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/mappers"
)

type UserRepository struct {
	datasource oracleDataSource.UserDataSource
}

func NewUserRepository(datasource oracleDataSource.UserDataSource) UserRepository {
	return UserRepository{datasource: datasource}
}

func (u UserRepository) SaveUserCompletely(ctx context.Context, user entities.User,basicDataUser entities.BasicUserData,userLocation entities.UserLocation,applicant entities.Applicant) (nis int, err error) {
	userModel := mappers.FromUserDomainToModel(user)
	//peopleModel := mappers.FromPeopleDomainToModel(people)
	basicDataUserModel := mappers.FromBasicDataUserDomainToModel(basicDataUser)
	userLocationModel := mappers.FromUserLocationDomainToModel(userLocation)
	applicantModel := mappers.FromApplicantDomainToModel(applicant)
	nis, err = u.datasource.SaveUserCompletely(ctx, userModel,basicDataUserModel,userLocationModel,applicantModel)
	if  err != nil {
		return nis, cockroachdbErrors.WithStack(err)
	}
	return nis, nil
}
