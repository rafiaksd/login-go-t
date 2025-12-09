package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("./users.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting DB: %v", err)
	}

	// Optional: configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
}
