package server

import (
	"fmt"
	"net/http"

	_"01.kood.tech/git/kretesaak/forum/internal/database"
	"01.kood.tech/git/kretesaak/forum/internal/registration"
)

func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")

	// Error handling with wrong path
	if r.URL.Path != "/loginauth" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
		return
	}

	// Errpr handling wrong method
	if r.Method != "POST" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	// Username criteria
	username := r.FormValue("usernameIn")
	fmt.Println("User logged in:", username)

	// Password criteria
	password := r.FormValue("passwordIn")
	fmt.Println("Password used by user:", password)

	// TODO checki kas sellinne username ja pass on baasis (hash) olemas ja saada ta õige puhul landing pagele oma eriliste kasutajaomadustega

	// login connection
	err := tmpl.ExecuteTemplate(w, "loginauth", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerAuthHandler running*****")

	// Error handling with wrong path
	if r.URL.Path != "/registerauth" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
		return
	}

	// Errpr handling wrong method
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

	// login connection
	err := tmpl.ExecuteTemplate(w, "registerauth", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
	var rf database.User

	var uID string
	rf.Id = uID // TODO vist?

	rf.Name = username
	rf.Email = email
	rf.Password = password

	fmt.Println("-----")
	fmt.Println(rf)
	fmt.Println("-----")

	// TODO checkid vastu baasi
	rslt := database.DbGetUserByIdOrEmail(rf.Email)
	fmt.Println(rslt)
	fmt.Println("*****")

	// Hashing test
	hpass := database.HashPassword(rf.Password)
	fmt.Println("Hash password:", hpass)
	fmt.Println("Hash correct:", database.CheckPasswordHash(rf.Password, hpass))
	*/

	/*

		// login connection
		err := tmpl.ExecuteTemplate(w, "registerauth", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	// TODO kui vorm on norm, siis ta peaks saatma registerAuth lehele vms, mis ütleb et kõik on norm ja kus on nupp mis suunab logimise lehele tagasi
	// TODO Või siis kohe login page-le, kus üleval on template teade, et account created succesfully.

}
