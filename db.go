package main

import (
	"fmt"
	"log"

	"github.com/zilstream/bot/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase() *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		helpers.GetEnv("DB_HOST", "localhost"),
		helpers.GetEnv("DB_USER", "melvin"),
		helpers.GetEnv("DB_PASSWORD", "test"),
		helpers.GetEnv("DB_NAME", "tinyzil"),
		helpers.GetEnv("DB_PORT", "5432"),
		helpers.GetEnv("DB_SSLMODE", "disable"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	runMigrations(db)

	return db
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate()
}
