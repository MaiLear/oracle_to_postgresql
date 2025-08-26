package config

import (
	"log"
	"github.com/joho/godotenv"
	"gitlab.com/sofia-plus/oracle_to_postgresql/infrastructure/loggers"
)

func Load(){
	_ = godotenv.Load()
	// Inicializar logger
	if err := loggers.InitLogger("server"); err != nil {
		log.Fatalf("Error inicializando logger: %v", err)
	}
}
