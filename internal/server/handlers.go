package server

import (
	"fmt"
	"net/http"

	"01.kood.tech/git/kretesaak/forum/internal/registration"
)

// Registerhandler serves a from for registering new users
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tmpl.ExecuteTemplate(w, "register", nil)
}

// registerAuthHandler creates new user in database
func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerAuthHandler running*****")
	// Getting form data
	r.ParseForm()

	username := r.FormValue("usernameUp")
	fmt.Println(username)
	// TODO päras on if shortcircuit?
	ua := registration.UsernameCorrect(username)
	fmt.Println("Username alphanumeric: ", ua)

	ul := registration.UsernameLen(username)
	fmt.Println("Username length: ", ul)

	// check password criteria
	password := r.FormValue("passwordUp")
	fmt.Println(password)
	pwc := registration.PswdConditions(password)
	fmt.Println(pwc)
	fmt.Println("password:", password, "\npswdLength:", len(password))

	// Username has missing criteria
	if !ua || !ul {
		fmt.Println("Username has missing criteria")
		// TODO see tekst peaks ilmuma normaalsesse kohta
		tmpl.ExecuteTemplate(w, "register", "Please check username criteria")
	}

	// Password criteria
	if !pwc.Lowercase || !pwc.Uppercase || !pwc.Number || !pwc.Special || !pwc.Length || pwc.NoSpaces {
		fmt.Println("Password has missing criteria")
		tmpl.ExecuteTemplate(w, "register", "Please check password criteria")
	}
}
