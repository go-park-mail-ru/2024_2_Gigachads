package httpserver

import (
	"net/http"
	"net/mail"
	"regexp"
	//"fmt"
)

type User struct {
    login string
    name string
    password string
}

var UserDB = make(map[string]User)

var signupFormTmpl = []byte(`
<html>
	<body>
	<form action="/signup" method="post">
		Login: <input type="text" name="login">
		Name: <input type="text" name="name">
		Password: <input type="password" name="password">
		Repeat Password: <input type="password" name="repassword">
		<input type="submit" value="Sign Up">
	</form>
	</body>
</html>
`) //TODO: убрать когда будет фронт регистрации

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(signupFormTmpl) //TODO: убрать когда будет фронт регистрации
		return
	}

	r.ParseForm()

	inputLogin := r.FormValue("login")
	inputName := r.FormValue("name")
	inputPassword := r.FormValue("password")
	inputRePassword := r.FormValue("repassword")

	if !emailIsValid(inputLogin) {
		ErrorResponse(w, r, "invalid_email")
		return
	}

	if !inputIsValid(inputName) {
		ErrorResponse(w, r, "invalid_input")
		return
	}

	if !inputIsValid(inputPassword) {
		ErrorResponse(w, r, "invalid_input")
		return
	}// а нужно ли?

	if !inputIsValid(inputRePassword) {
		ErrorResponse(w, r, "invalid_input")
		return
	}// а нужно ли?

	if inputPassword != inputRePassword {
		ErrorResponse(w, r, "invalid_password")
		return
	}

	if _, ok := UserDB[inputLogin]; ok {
		ErrorResponse(w, r, "login_taken")
		return
	}

	UserDB[inputLogin] = User{login: inputLogin, name: inputName, password: inputPassword}
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte(inputLogin))
	//fmt.Fprintln(w, UserDB)
	
}

func emailIsValid(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}

func inputIsValid(str string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", str)
	return match
}