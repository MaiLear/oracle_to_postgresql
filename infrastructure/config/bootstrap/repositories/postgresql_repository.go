package repositories

import (
	"log"

	dbConfig "gitlab.com/sofia-plus/go_db_connectors/config"
	"gitlab.com/sofia-plus/oracle_to_postgresql/domain/ports/repositories"
	"gitlab.com/sofia-plus/oracle_to_postgresql/interface_adapters/gateways/db/repositories/postgresql"
	"gorm.io/gorm"
)

type PostgresqlRepository struct {
	TrainingProgram repositories.RepositoryWrite
}

func initPosConnection() (connection *gorm.DB) {
	connection, err := dbConfig.NewPostgresConnection()
	if err != nil {
		log.Fatalf("infra: fallo al conectar a PostgreSQL -> %v", err)
	}
	return connection
}

func InitPosRepository() PostgresqlRepository {
	dbConnection := initPosConnection()
	trainingProgram := postgresql.NewTrainingProgram(dbConnection)
	return PostgresqlRepository{
		TrainingProgram: trainingProgram,
	}
}
