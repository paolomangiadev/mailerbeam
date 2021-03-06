package models

import (
	"log"

	uuid "github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name     string    `gorm:"not null;default: null" json:"name"`
	Username string    `gorm:"unique;not null;default: null" json:"username"`
	Password string    `gorm:"type:varchar(255);not null;default: null" json:"password"`
	Email    string    `gorm:"type:varchar(100);unique_index;default: null" json:"email"`
	Role     string    `gorm:"not null;default: null" json:"role"`
}

// CreateUserRequest type
type CreateUserRequest struct {
	Email    string `json:"email" valid:"required,email"`
	Name     string `json:"name" valid:"required,stringlength(2|50)"`
	Username string `json:"username" valid:"required,stringlength(2|50)"`
	Password string `json:"password" valid:"required,stringlength(8|50)"`
}

// LogingUserRequest type
type LogingUserRequest struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required"`
}

// hash password with bcrypt
func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, 10)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// ComparePasswords with plain pw
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// BeforeSave func
func (u *User) BeforeSave() (err error) {
	hash := hashAndSalt([]byte(u.Password))
	u.Password = hash
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	u1, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	scope.SetColumn("ID", u1)
	return nil
}
