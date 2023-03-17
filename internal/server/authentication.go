package server

import (
	"fmt"
	"net/http"

	"01.kood.tech/git/kretesaak/forum/internal/registration"
)

type RegistrationForm struct {
	FormUsername string
	FormEmail    string
	FormPassword string
}

func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerAuthHandler running*****")

	// Error handling with wrong path
	if r.URL.Path != "/registerauth" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	// Username criteria
	username := r.FormValue("usernameUp")
	ua := registration.UsernameCorrect(username)
	ul := registration.UsernameLen(username)

	// Email criteria
	email := r.FormValue("emailUp")
	ev := registration.IsValidEmail(email)

	// Password criteria
	password := r.FormValue("passwordUp")
	pwc := registration.PswdConditions(password)

	// Username has missing criteria
	if !ua || !ul {
		fmt.Println("Username has missing criteria")
		// TODO see tekst peaks ilmuma normaalsesse kohta
		tmpl.ExecuteTemplate(w, "register", "Please check username criteria")
	}

	// Email criteria
	if !ev {
		fmt.Println("Email has missing criteria")
		tmpl.ExecuteTemplate(w, "register", "Please check email")
	}

	// Password criteria
	if !pwc.Lowercase || !pwc.Uppercase || !pwc.Number || !pwc.Special || !pwc.Length || pwc.NoSpaces {
		fmt.Println("Password has missing criteria")
		tmpl.ExecuteTemplate(w, "register", "Please check password criteria")
	}

	var rf RegistrationForm

	rf.FormUsername = username
	rf.FormEmail = email
	rf.FormPassword = password

	fmt.Println("-----")
	fmt.Println(rf)
	fmt.Println("-----")

	// TODO checkid vastu baasi
	


	// TODO kui vorm on norm, siis ta peaks saatma registerAuth lehele vms, mis ütleb et kõik on norm ja kus on nupp mis suunab logimise lehele tagasi
	// TODO Või siis kohe login page-le, kus üleval on template teade, et account created succesfully.

}

