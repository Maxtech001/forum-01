package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"01.kood.tech/git/kretesaak/forum/internal/database"
	_ "github.com/gofrs/uuid"
)

func getUserByCookie(r *http.Request) string {
	result := ""
	cookie, err := r.Cookie("session")
	if err != nil {
		return result
	}
	result = database.DbGetUserByCookie(cookie.Value)
	return result
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	user_id := getUserByCookie(r)

	// Error handling with wrong path
	if r.URL.Path != "/" {
		tmpl.ExecuteTemplate(w, "error", user_id)
		return
	}
	// Wrong method handling
	if r.Method != "GET" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Getting content
	mainPageContent := getMainPageContent(user_id, r.URL.Query())
	/*
		var mainPageContent database.Mainpage
		mainPageContent.User_id = user_id
		mainPageContent.Posts = database.DbGetPosts()
		mainPageContent.Tags = database.DbGetTags()
	*/

	err := tmpl.ExecuteTemplate(w, "index", mainPageContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	user_id := getUserByCookie(r)
	if user_id == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Wrong method handling
	if r.Method != "GET" {
		http.Error(w, "Bad request - method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Extract post ID from URL
	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/"))
	if err != nil {
		http.Error(w, "Bad request - invalid post ID.", http.StatusBadRequest)
		return
	}

	// Getting content
	err, postPageContent := getPostPageContent(postID, user_id)

	// If missing
	if err != nil {
		tmpl.ExecuteTemplate(w, "error", user_id)
		fmt.Println(http.StatusNotFound)
		return
	}

	// Render post template
	err = tmpl.ExecuteTemplate(w, "comment", postPageContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	user_id := getUserByCookie(r)

	// Wrong method handling
	if r.Method != "GET" {
		http.Error(w, "Bad request - method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Extract post ID from URL
	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		http.Error(w, "Bad request - invalid post ID.", http.StatusBadRequest)
		return
	}

	// Getting content
	err, postPageContent := getPostPageContent(postID, user_id)

	// If missing
	if err != nil {
		tmpl.ExecuteTemplate(w, "error", user_id)
		fmt.Println(http.StatusNotFound)
		return
	}

	// Render post template
	err = tmpl.ExecuteTemplate(w, "post", postPageContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	user_id := getUserByCookie(r)

	if r.Method != "POST" {
		// Error handling with wrong path
		if r.URL.Path != "/createpost" {
			tmpl.ExecuteTemplate(w, "error", user_id)
			fmt.Println(http.StatusNotFound)
			return
		}

		// Getting content
		createPostPageContent := getCreatePostPageContent(user_id)

		// login connection
		err := tmpl.ExecuteTemplate(w, "createpost", createPostPageContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user_id := getUserByCookie(r)

	// Error handling with wrong path
	if r.URL.Path != "/login" {
		tmpl.ExecuteTemplate(w, "error", user_id)
		fmt.Println(http.StatusNotFound)
		return
	}
	// Wrong method handling
	if r.Method != "GET" {
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

func aboutUsHandler(w http.ResponseWriter, r *http.Request) {
	user_id := getUserByCookie(r)
	readme, err := os.ReadFile("README.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "about", database.Aboutpage{User_id: user_id, Content: string(readme)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// TODO ta teeb praegu miskipärast kaks korda seda päringut (teine on tühi), luua login ja reg endpoint erinevalt esialgu
// registerAuthHandler creates new user in database
func registerHandler(w http.ResponseWriter, r *http.Request) {
	user_id := getUserByCookie(r)

	// Error handling with wrong path
	if r.URL.Path != "/register" {
		tmpl.ExecuteTemplate(w, "error", user_id)
		fmt.Println(http.StatusNotFound)
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
