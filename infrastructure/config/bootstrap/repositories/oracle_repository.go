package repositories

import (
	"log"

	dbConfig "gitlab.com/sofia-plus/go_db_connectors/config"
	"gitlab.com/sofia-plus/oracle_to_postgresql/domain/ports/repositories"
	"gitlab.com/sofia-plus/oracle_to_postgresql/interface_adapters/gateways/db/repositories/oracle"
	"gorm.io/gorm"
)

type OracleRepository struct {
	TrainingProgram repositories.RepositoryReader
}

func initOraConnection() (connection *gorm.DB) {
	connection, err := dbConfig.NewOracleConnection()
	if err != nil {
		log.Fatalf("infra: fallo al conectar a Oracle -> %v", err)
	}
	return connection
}
func InitOraRepository() OracleRepository {
	dbConnection := initOraConnection()
	trainingProgram := oracle.NewTrainingProgram(dbConnection)
	return OracleRepository{
		TrainingProgram: trainingProgram,
	}
}
