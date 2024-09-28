package database

type User struct {
	Id string
	Name     string 
	Email    string 
	Password string
}

var UserDB = make(map[string]User) //найти айди по юзеру

var UserID = make(map[string]string) //найти юзера по айди