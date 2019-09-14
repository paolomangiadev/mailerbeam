package models

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

var db *gorm.DB

func Init() {
	dbName := os.Getenv("DB_NAME")

	conn, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		panic("failed to connect database")
	}
	db = conn

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create persons table
		{
			ID: "1",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists(&User{}, "users").Error
			},
		},
	})

	errMigrate := m.Migrate()
	if errMigrate != nil {
		log.Printf("Could not migrate: %v", err)
	} else {
		log.Printf("Migration did run successfully")
	}
	defer conn.Close()
}

func GetDB() *gorm.DB {
	return db
}
