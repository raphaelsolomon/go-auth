package models

type UserResponse interface {
	GetUser() *User
}

type User struct {
	ID              uint   `gorm:"prinmaryKey"`
	Uuid            string `gorm:"unique"`
	Email           string `gorm:"unique"`
	Password        string `json:"-"`
	Firstname       string
	Lastname        string
	Status          bool
	PhoneNo         string `gorm:"unique"`
	IsEmailverified bool
}

func (u *User) GetUser() *User {
	return &User{
		ID:              u.ID,
		Uuid:            u.Uuid,
		Email:           u.Email,
		Firstname:       u.Firstname,
		Lastname:        u.Lastname,
		Status:          u.Status,
		PhoneNo:         u.PhoneNo,
		IsEmailverified: u.IsEmailverified,
	}
}
