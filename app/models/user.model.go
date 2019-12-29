package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// User model
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name     string    `gorm:"not null;default: null" json:"name"`
	Username string    `gorm:"unique;not null;default: null" json:"username"`
	Password string    `gorm:"not null;default: null" json:"password"`
	Email    string    `gorm:"type:varchar(100);unique_index;default: null" json:"email"`
	Role     string    `gorm:"size:255;not null;default: null" json:"role"`
}

// CreateUserRequest type
type CreateUserRequest struct {
	Email    string `json:"email" valid:"required,email"`
	Name     string `json:"name" valid:"required,stringlength(2|50)"`
	Username string `json:"username" valid:"required,stringlength(2|50)"`
	Password string `json:"password" valid:"required,stringlength(8|50)"`
}

// BeforeSave func
func (u *User) BeforeSave() (err error) {
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	return nil
}
