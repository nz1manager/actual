package config

import (
    "log"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "yourmodule/internal/models"
)

func InitDB() *gorm.DB {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL environment variable not set")
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate schemas
    if err := db.AutoMigrate(&models.User{}, &models.Test{}, &models.Submission{}); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    return db
}
