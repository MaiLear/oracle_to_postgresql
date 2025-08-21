package services

import (
	"flag"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/application/services"
	inputPort "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/in"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/config/bootstrap/repositories"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/services/email"
)

type Service struct {
	MainService inputPort.ServiceChain
}

func getConsoleVariables() map[string]any {
	updateCharacterizationSheet := flag.Bool("update", false, "Actualizar ficha de caracterizacion")
	flag.Parse()
	return map[string]any{
		"update": *updateCharacterizationSheet,
	}
}

func InitService() Service {
	postgresqlRepo := repositories.InitPosRepository()
	oracleRepo := repositories.InitOraRepository()
	consoleVariables := getConsoleVariables()
	errorService := services.NewErrorService(postgresqlRepo.ErrorRepository)
	peopleService := services.NewPeopleService(postgresqlRepo.PeopleRepository, oracleRepo.PeopleRepository, errorService)
	applicantService := services.NewApplicantService(oracleRepo.ApplicantRepository, postgresqlRepo.ApplicantRepository, errorService)
	emailService, err := email.NewSmtpService()
	if err != nil {
		panic("Error inicializando servicio de email: " + err.Error())
	}
	basicDataUserService := services.NewBasicDataUserService(postgresqlRepo.BasicDataUserRepository, errorService)
	userLocationService := services.NewUserLocationService(postgresqlRepo.UserLocationRepository, errorService)
	userService := services.NewUserService(postgresqlRepo.UserRepository, oracleRepo.UserRepository, userLocationService, peopleService, basicDataUserService, applicantService, errorService)
	parameterService := services.NewParameterService(
		oracleRepo.ParameterRepository,
		postgresqlRepo.ParameterRepository,
		oracleRepo.DocumentTypeRepository,
		postgresqlRepo.DocumentTypeRepository,
	)

	enrollmentValidationService := services.NewEnrollmentValidationService(
		oracleRepo.EnrollmentValidationRepository,
		parameterService,
	)

	enrollmentService := services.NewEnrollmentService(postgresqlRepo.EnrollmentRepository, oracleRepo.EnrollmentRepository, peopleService, applicantService, emailService, enrollmentValidationService,errorService, consoleVariables["update"].(bool))
	mainService := services.NewService([]inputPort.ServiceChain{
		userService,
		enrollmentService,
	})
	return Service{
		MainService: mainService,
	}
}
