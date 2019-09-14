package models

type User struct {
	Id       string `gorm:"type:uuid;primary_key;"`
	Name     string
	Username string
	Password string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Role     string `gorm:"size:255"`
}

func (u *User) BeforeSave() (err error) {
	return
}
