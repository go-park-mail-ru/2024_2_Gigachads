package models

type User struct {
	ID         int
	Email      string `json:"email"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
	AvatarURL  string `json:"avatarPath"`
}

type UserRepository interface {
	CreateUser(signup *User) (*User, error)
	CheckUser(login *User) (*User, error)
	IsExist(email string) (bool, error)
	UpdateInfo(user *User) error
	GetUserByEmail(email string) (*User, error)
}
