package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string
	Username string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Role     string `gorm:"size:255"`
}
