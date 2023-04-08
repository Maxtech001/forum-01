package server

import (
	"fmt"
	"net/http"

	"01.kood.tech/git/kretesaak/forum/internal/database"
	"01.kood.tech/git/kretesaak/forum/internal/registration"
	"github.com/gofrs/uuid"
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
	email := r.FormValue("emailIn")
	fmt.Println("User logged in:", email)

	// Password criteria
	password := r.FormValue("passwordIn")
	fmt.Println("Password used by user:", password)

	cp, user_id := database.DbAuthenticateUser(email, password)

	// TODO või siis email või password on vale indikaator saata kasutajale
	if !cp {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	cookie, err := r.Cookie("session")
	fmt.Println("Cookie get:", cookie)
	if err != nil {
		// Creating a version 4 UUID
		id, err2 := uuid.NewV4()
		if err2 != nil {
			fmt.Printf("failed to generate UUID: %v", err2)
		}
		fmt.Printf("generated Version 4 UUID %v", id)

		// TODO expiration
		cookie = &http.Cookie{
			Name:  "session",
			Value: id.String(),
			// Secure: true
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie) // setting a cookie if it does not exist
		fmt.Println("Cookie set:", cookie)
		database.DbAddCookie(cookie.Value, user_id)
	}

	// login connection
	err = tmpl.ExecuteTemplate(w, "loginauth", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****logout running*****")

	// Error handling with wrong path
	if r.URL.Path != "/logout" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
		return
	}

	// Errpr handling wrong method
	if r.Method != "GET" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	cookie, _ := r.Cookie("session")
	database.DbDeleteCookie(cookie.Value)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	// login connection
	err := tmpl.ExecuteTemplate(w, "logout", nil)
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
	/*
		Front-end lookup
	*/
	r.ParseForm()
	// Email, username and password
	username, email, password := r.FormValue("usernameUp"), r.FormValue("emailUp"), r.FormValue("passwordUp")
	/*
		Database lookup
	*/
	rf := registration.NewUser(username, email, password)

	// Checking email and username
	ux := database.DbUserIdExist(rf.Id)
	fmt.Println("Is username:", ux)
	ex := database.DbEmailExist(rf.Email)
	fmt.Println("Is email:", ex)
	// TODO see indikaator saata kasutajale
	if ux || ex {
		fmt.Println("Email or username already exists")
		tmpl.ExecuteTemplate(w, "register", "Please choose another email and/or username because it is already registered")
		return
	}
	// Inserting values into database
	fmt.Println(database.DbInsertUser(rf))
	// Going to login page
	err := tmpl.ExecuteTemplate(w, "registerauth", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
