package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	email, password := r.FormValue("emailIn"), r.FormValue("passwordIn")
	fmt.Println("User logged in:", email)
	fmt.Println("Password used by user:", password)

	cp, user_id := database.DbAuthenticateUser(email, password)

	// TODO või siis email või password on vale indikaator saata kasutajale
	if !cp {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	/*
		Cookie logic
	*/
	cookie, err := r.Cookie("session")
	fmt.Println("Cookie get:", cookie)
	if err != nil {
		// Creating a version 4 UUID
		id, err2 := uuid.NewV4()
		if err2 != nil {
			fmt.Printf("failed to generate UUID: %v", err2)
		}
		exp := time.Now().Add(24 * time.Hour)
		fmt.Println(exp)
		exp = exp.UTC()
		fmt.Printf("generated Version 4 UUID %v\nexpires UTC: %v\n", id, exp.Format("2006-01-02 15:04:05"))

		cookie = &http.Cookie{
			Name:  "session",
			Value: id.String(),
			// Secure: true
			HttpOnly: true,
			Path:     "/",
			Expires:  exp,
		}
		http.SetCookie(w, cookie) // setting a cookie if it does not exist
		//	fmt.Println("Cookie set:", cookie)
		database.DbAddCookie(cookie.Value, user_id, exp)
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
		tmpl.ExecuteTemplate(w, "register", "Please choose another email and/or username because it already exists")
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

func commentAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Errpr handling wrong method
	if r.Method != "POST" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/commentauth/"))
	if err != nil {
		http.Error(w, "Bad request - invalid post ID.", http.StatusBadRequest)
		return
	}

	r.ParseForm()

	err = tmpl.ExecuteTemplate(w, "commentauth", postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
