package server

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

// TODO andmebaasi konvertimine
var dbSessions = map[string]string{}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")

	// Error handling with wrong path
	if r.URL.Path != "/login" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
		return
	}
	// Wrong method handling
	if r.Method != "GET"{
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// login connection
	err := tmpl.ExecuteTemplate(w, "login", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TODO ta teeb praegu miskipärast kaks korda seda päringut (teine on tühi), luua login ja reg endpoint erinevalt esialgu
// registerAuthHandler creates new user in database
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")

	// Error handling with wrong path
	if r.URL.Path != "/register" {

		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
		return
	}
	// Wrong method handling
	if r.Method != "GET" {

		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// TODO ilmselt kuidagi username põhiselt pärast logimise tegemist
	////////////////////// Cookie generation
	cookie, err := r.Cookie("session")

	if err != nil {
		// Creating a version 4 UUID
		id, err2 := uuid.NewV4()
		if err2 != nil {
			fmt.Printf("failed to generate UUID: %v", err2)
		}
		fmt.Printf("generated Version 4 UUID %v", id)
		
		// TODO expiration, logout, clearing, specific sites (regamis ja login saitidel ilmselt kaob ära)
		cookie = &http.Cookie{
			Name: "session",
			Value: id.String(),
			// Secure: true
			HttpOnly: true,
			Path: "/",
		}
		http.SetCookie(w, cookie) // setting a cookie if it does not exist
	}
	fmt.Println(cookie)
	// TODO andmebaasi viia ja parssimine korda teha
	// Getting the user and assigning cookie
	/*
	// if the user exists already, get the user
	var u user
	if un, ok := dbSessions[cookie.Value]; ok {
		u = dbUsers[un] // user_id saab ülejäänud info kasutaja kohta kätte
	}

	// Pärast logimist või authentication lehel määrata
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		f := r.FormValue("firstname")
		l := r.FormValue("lastname")
		u = user{un, f, l}
		dbSessions[cookie.Value] = un
		dbUsers[un] = u
	}
	*/
	// TODO ja siis template execution põhineb kasutaja structil, ilmselt kui puudub siis nil või muu leht
	//err = tmpl.ExecuteTemplate(w, "index", u)
	////////////////////// Cookie generation



	// Register connection
	err = tmpl.ExecuteTemplate(w, "register", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Ilmselt sisselogimisel klõpsamine
// Another page redirection, cookie otsimine
/*
func bar(w http.ResponseWriter, r *http.Request) {

	// get cookie
	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// TODO siin on ilmselt mingi timeout
	un, ok := dbSessions[c.Value]
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u := dbUsers[un]
	tmpl.ExecuteTemplate(w, "bar", u)

}
*/
