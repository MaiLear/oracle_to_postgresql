package main

import (
	"log"

	"github.com/joho/godotenv"
	serviceConfiguration "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/config/bootstrap/services"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/logger"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/presentation/handlers"
)

func main() {
	_ = godotenv.Load()
	// Inicializar logger
	if err := logger.InitLogger("server"); err != nil {
		log.Fatalf("Error inicializando logger: %v", err)
	}
	services := serviceConfiguration.InitService()
	mainHandler := handlers.NewHandler(services.MainService)
	mainHandler.Execute()
}
