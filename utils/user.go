package utils

// Abstract interface for a User
type IUser interface {
	GetUser() *User
	SetFirstname(value string)
	SetLasttname(value string)
	SetEmail(value string)
	SetPassword(value string)
}

type User struct {
	firstname string
	lastname  string
	email     string
	password  string
	Age       int
}

func NewUser(firstname, lastname, email, password string) *User {
	return &User{firstname: firstname, lastname: lastname, email: email, password: password}
}

func (u *User) GetUser() *User {
	return u
}

func (u *User) SetFirstname(value string) {
	u.firstname = value
}

func (u *User) SetLasttname(value string) {
	u.lastname = value
}

func (u *User) SetEmail(value string) {
	u.email = value
}

func (u *User) SetPassword(value string) {
	u.password = value
}
