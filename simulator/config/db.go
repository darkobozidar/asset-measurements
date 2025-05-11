package config  // TODO rename to init

import (
	"fmt"
	"log"

	"gorm.io/gorm"
    "gorm.io/driver/postgres"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	// TODO collect these from .env.
	// TODO put to constant the base sting.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		"postgres",  // TODO check if this is correct
		"admin",
		"admin",
		"asset_measurements",
		"5432",
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Connecting to DB failed.", dsn)
	}
}
