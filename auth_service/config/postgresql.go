package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresqlDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Jakarta", Envs.POSTGRESQL_HOST, Envs.POSTGRESQL_USER, Envs.POSTGRESQL_PASSWORD, Envs.POSTGRESQL_DB, Envs.POSTGRESQL_PORT)
	var err error
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	return DB
}
