package models

// User model
type User struct {
	ID       string `gorm:"type:uuid;primary_key;"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `gorm:"type:varchar(100);unique_index" json:"email"`
	Role     string `gorm:"size:255" json:"role"`
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
