package models

import (
	"os"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Init db
func Init() {
	dbName := os.Getenv("DB_NAME")

	conn, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		panic("failed to connect database")
	}
	db = conn
}

// GetDB connection
func GetDB() *gorm.DB {
	return db
}

// CloseDB connection
func CloseDB() error {
	return db.Close()
}
