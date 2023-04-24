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

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	userId := getUserByCookie(r)

	// Error handling with wrong path
	if r.URL.Path != "/" {
		tmpl.ExecuteTemplate(w, "error", userId)
		return
	}
	// Wrong method handling
	if r.Method != "GET" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Getting content
	mainPageContent := getMainPageContent(userId, r.URL.Query())

	err := tmpl.ExecuteTemplate(w, "index", mainPageContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	userId := getUserByCookie(r)

	// Redirect unregistered users
	if userId == "" {
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
	err, postPageContent := getPostPageContent(postID, userId)

	// If missing
	if err != nil {
		tmpl.ExecuteTemplate(w, "error", userId)
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
	userId := getUserByCookie(r)

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
	err, postPageContent := getPostPageContent(postID, userId)

	// If missing
	if err != nil {
		tmpl.ExecuteTemplate(w, "error", userId)
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
	userId := getUserByCookie(r)

	// Redirect unregistered users
	if userId == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		// Error handling with wrong path
		if r.URL.Path != "/createpost" {
			tmpl.ExecuteTemplate(w, "error", userId)
			fmt.Println(http.StatusNotFound)
			return
		}

		// Getting content
		createPostPageContent := getCreatePostPageContent(userId)

		// Create post connection
		err := tmpl.ExecuteTemplate(w, "createpost", createPostPageContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	userId := getUserByCookie(r)

	// Redirect if already logged in
	if userId != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Error handling with wrong path
	if r.URL.Path != "/login" {
		tmpl.ExecuteTemplate(w, "error", userId)
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
	userId := getUserByCookie(r)
	readme, err := os.ReadFile("README.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "about", database.Aboutpage{User_id: userId, Content: string(readme)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	userId := getUserByCookie(r)

	// Redirect if already logged in user
	if userId != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Error handling with wrong path
	if r.URL.Path != "/register" {
		tmpl.ExecuteTemplate(w, "error", userId)
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
