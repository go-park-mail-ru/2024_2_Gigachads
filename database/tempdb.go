package database

type User struct {
	Name     string 
	Email    string 
	Password string
}

var UserDB = make(map[string]User) //найти user по email

var UserHash = make(map[string]string) //найти email gо хэшу