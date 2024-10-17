package models

type User struct {
	Name     string
	Email    string
	Password string
}

type Signup struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
