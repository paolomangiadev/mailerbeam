package models

import (
	"os"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func Init() {
	dbName := os.Getenv("DB_NAME")

	conn, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		panic("failed to connect database")
	}
	db = conn
	defer conn.Close()
}

func GetDB() *gorm.DB {
	return db
}
