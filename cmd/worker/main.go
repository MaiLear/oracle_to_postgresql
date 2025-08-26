package main

import (
	"gitlab.com/sofia-plus/oracle_to_postgresql/infrastructure/config"
	serviceConfiguration "gitlab.com/sofia-plus/oracle_to_postgresql/infrastructure/config/bootstrap/services"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/presentation/handlers"
)

func main() {
	config.Load()
	services := serviceConfiguration.InitService()
	controller := handlers.NewHandler(services.MainService)
}
