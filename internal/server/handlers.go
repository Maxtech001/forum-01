package server

import (
	"fmt"
	"net/http"
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
	fmt.Println("getUserByCookie DB request", cookie.Value)
	result = database.DbGetUserByCookie(cookie.Value)
	fmt.Println("getUserByCookie result:", result)
	return result
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	user_id := getUserByCookie(r)

	// Error handling with wrong path
	if r.URL.Path != "/" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
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
		http.Error(w, "Bad request - post not found.", http.StatusNotFound)
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
	fmt.Println("*****createPostHandler running*****")
	user_id := getUserByCookie(r)

	if r.Method != "POST" {
		// Error handling with wrong path
		if r.URL.Path != "/createpost" {
			http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
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
	} else {
		user_id := getUserByCookie(r)
		title := r.FormValue("titleIn")
		content := r.FormValue("contentIn")
		r.ParseForm()
		tags2 := r.Form["tag"]
		var tags1 []int
		for _, i := range tags2 {
			j, err := strconv.Atoi(i)
			if err != nil {
				panic(err)
			}
			tags1 = append(tags1, j)
		}

		// Validate form data
		if title == "" || content == "" {
			http.Error(w, "Please fill in all fields", http.StatusBadRequest)
			return
		}

		err, post_id := database.DbInsertPost(user_id, title, content, tags1)
		if err != nil {
			fmt.Println("DbInsertpost Error")
		} else {
			fmt.Println("DbInsertpost Success:", post_id)
		}

		// Redirect to success page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")

	// Error handling with wrong path
	if r.URL.Path != "/login" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
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
