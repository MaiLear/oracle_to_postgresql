package usecases

import (
	"flag"

	"gitlab.com/sofia-plus/oracle_to_postgresql/infrastructure/config/bootstrap/repositories"
	"gitlab.com/sofia-plus/oracle_to_postgresql/usecases"
	inputPort "gitlab.com/sofia-plus/oracle_to_postgresql/usecases/ports/in"
	"gitlab.com/sofia-plus/oracle_to_postgresql/usecases/ports/pipeline"
)

type Service struct {
	Usecase inputPort.Port
}

func getConsoleVariables() map[string]any {
	updateCharacterizationSheet := flag.Bool("update", false, "Actualizar ficha de caracterizacion")
	flag.Parse()
	return map[string]any{
		"update": *updateCharacterizationSheet,
	}
}

func InitUsecases() Service {
	postgresqlRepo := repositories.InitPosRepository()
	oracleRepo := repositories.InitOraRepository()
	traningProgram := usecases.NewTrainingProgram(oracleRepo.TrainingProgram,postgresqlRepo.TrainingProgram)
	mainUseCase := usecases.NewUseCase([]pipeline.Service{
		traningProgram,
	} )

	return Service{
		Usecase: mainUseCase,
	}
}
