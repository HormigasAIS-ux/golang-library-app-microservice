package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresqlDB() *gorm.DB {
	// create instance of default database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Jakarta",
		Envs.POSTGRESQL_HOST,
		Envs.POSTGRESQL_USER,
		Envs.POSTGRESQL_PASSWORD,
		"postgres",
		Envs.POSTGRESQL_PORT,
	)

	logger.Debugf("connecting to default database: postgres")
	var err error
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("failed to connect to the database: %v", err)
	}

	// create database from environment variable to ensure it exists
	logger.Debugf("ensuring database: %s", Envs.POSTGRESQL_DB)
	err = DB.Exec(fmt.Sprintf("CREATE DATABASE %s", Envs.POSTGRESQL_DB)).Error
	if err != nil {
		logger.Warningf("Failed to create database: %v", err)
	}

	// recreate database using database name from environment variable
	dsn = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Jakarta",
		Envs.POSTGRESQL_HOST,
		Envs.POSTGRESQL_USER,
		Envs.POSTGRESQL_PASSWORD,
		Envs.POSTGRESQL_DB,
		Envs.POSTGRESQL_PORT,
	)

	logger.Debugf("connecting to database: %s", dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("failed to connect to the database: %v", err)
	}

	return DB
}
