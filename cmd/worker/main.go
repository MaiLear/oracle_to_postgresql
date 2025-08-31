package main

import (
	"gitlab.com/sofia-plus/oracle_to_postgresql/infrastructure/config"
	usecaseConfiguration "gitlab.com/sofia-plus/oracle_to_postgresql/infrastructure/config/bootstrap/usecases"
	"gitlab.com/sofia-plus/oracle_to_postgresql/interface_adapters/controllers"
)

func main() {
	config.Load()
	usecases := usecaseConfiguration.InitUsecases()
	controller := controllers.NewController(usecases.Usecase)
	controller.Execute()
}
