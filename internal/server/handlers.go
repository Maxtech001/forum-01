package server

import (
	"fmt"
	"net/http"
)

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

	// Register connection
	err := tmpl.ExecuteTemplate(w, "register", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
