package config

import (
    "fmt"
    "log"
    "os"
    "database/sql"

    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

var DB *gorm.DB
var SQLDB *sql.DB

func ConnectToPostgresDB() {
    var err error
    const BASE_CONN_URL string = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"

    dsn := fmt.Sprintf(
        BASE_CONN_URL,
        os.Getenv("POSTGRES_HOST"),
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRES_DB"),
        os.Getenv("POSTGRES_PORT"),
    )

    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        log.Fatal("Connecting to DB failed.", dsn)
    }

    SQLDB, err = DB.DB()
    if err != nil {
        log.Fatal("Failed to get underlying DB from GORM:", err)
    }
}
