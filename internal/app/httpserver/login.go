package httpserver

/*
TODO:
-кука
*/

import (
	"net/http"
	"time"
	//"fmt"
)

var loginFormTmpl = []byte(`
<html>
	<body>
	<form action="/login" method="post">
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login">
	</form>
	</body>
</html>
`) //TODO: убрать когда будет фронт авторизации

func LogInHandler(w http.ResponseWriter, r *http.Request) {

	//UserDB["nick@giga-mail.ru"] = User{login: "nick@giga-mail.ru", name: "nick", password: "12345"} //убрать, когда будет бд
	
	if r.Method != http.MethodPost {
		w.Write(loginFormTmpl) //TODO: убрать когда будет фронт авторизации
		return
	}

	r.ParseForm()

	inputLogin := r.FormValue("login")
	inputPassword := r.FormValue("password")

	elem, ok := UserDB[inputLogin]

	if !emailIsValid(inputLogin) {
		ErrorResponse(w, r, "invalid_input")
		return
	}// а нужно ли?

	if !inputIsValid(inputPassword) {
		ErrorResponse(w, r, "invalid_input")
		return
	}// а нужно ли?

	if !ok {
		ErrorResponse(w, r, "user_does_not_exist")
		return
	}

	if UserDB[inputLogin].password != inputPassword {
		ErrorResponse(w, r, "invalid_password")
		return
	}

	w.WriteHeader(http.StatusOK)
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   elem.login,
		Expires: expiration,
		HttpOnly: true,
		//Domain: "127.0.0.1",
	}
	http.SetCookie(w, &cookie)
}