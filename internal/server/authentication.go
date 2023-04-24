package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"01.kood.tech/git/kretesaak/forum/internal/database"
	"github.com/gofrs/uuid"
)

func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Error handling with wrong path
	if r.URL.Path != "/loginauth" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
		return
	}

	// Error handling wrong method
	if r.Method != "POST" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	// Username criteria
	email, password := r.FormValue("emailIn"), r.FormValue("passwordIn")

	cp, user_id := database.DbAuthenticateUser(email, password)

	if !cp {
		fmt.Println("Email or password wrong")
		tmpl.ExecuteTemplate(w, "login", "Can't find a user with this email and password")
		return
	}
	/*
		Cookie logic
	*/
	cookie, err := r.Cookie("session")
	user := getUserByCookie(r)

	if err != nil || user == "" {
		// Creating a version 4 UUID
		id, err2 := uuid.NewV4()
		if err2 != nil {
			fmt.Printf("failed to generate UUID: %v", err2)
		}
		exp := time.Now().Add(24 * time.Hour)
		exp = exp.UTC()

		cookie = &http.Cookie{
			Name:  "session",
			Value: id.String(),
			// Secure: true
			HttpOnly: true,
			Path:     "/",
			Expires:  exp,
		}
		http.SetCookie(w, cookie) // setting a cookie if it does not exist
		database.DbAddCookie(cookie.Value, user_id, exp)
	}

	// Login connection
	err = tmpl.ExecuteTemplate(w, "loginauth", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Error handling with wrong path
	if r.URL.Path != "/logout" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
		return
	}

	// Error handling wrong method
	if r.Method != "GET" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Delete cookie
	cookie, _ := r.Cookie("session")
	database.DbDeleteCookie(cookie.Value)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	// Logout connection
	err := tmpl.ExecuteTemplate(w, "logout", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Error handling with wrong path
	if r.URL.Path != "/registerauth" {
		http.Error(w, "Bad request - 404 resource not found.", http.StatusNotFound)
		return
	}
	// Error handling wrong method
	if r.Method != "POST" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	/*
		Data varification in the frontend
	*/
	r.ParseForm()
	// Email, username and password
	username, email, password := r.FormValue("usernameUp"), r.FormValue("emailUp"), r.FormValue("passwordUp")
	/*
		Database lookup
	*/
	rf := newUser(username, email, password)

	// Checking if email and username already exist
	ux := database.DbUserIdExist(rf.Id)
	ex := database.DbEmailExist(rf.Email)

	if ux || ex {
		fmt.Println("Email or username already exists")
		tmpl.ExecuteTemplate(w, "register", "Please choose another email and/or username because it already exists")
		return
	}
	// Inserting values into database
	database.DbInsertUser(rf)
	// Going to login page
	err := tmpl.ExecuteTemplate(w, "registerauth", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func commentAuthHandler(w http.ResponseWriter, r *http.Request) {
	userId := getUserByCookie(r)

	// Redirect unregistered users
	if userId == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Errpr handling wrong method
	if r.Method != "POST" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Get post id from URL
	postId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/commentauth/"))
	if err != nil {
		http.Error(w, "Bad request - invalid post ID.", http.StatusBadRequest)
		return
	}

	// Insert comment to db
	content := r.FormValue("commentIn")
	err1 := database.DbInsertComment(postId, userId, content)
	if err1 != nil {
		fmt.Println("DbInsertComment Error")
	}

	// Send user to post page
	http.Redirect(w, r, "/post/"+strconv.Itoa(postId), http.StatusSeeOther)
}

func createPostAuthHandler(w http.ResponseWriter, r *http.Request) {
	userId := getUserByCookie(r)

	// Redirect unregistered users
	if userId == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Error handling wrong method
	if r.Method != "POST" {
		http.Error(w, "Bad request - 405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Add post to db
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

	err, _ := database.DbInsertPost(userId, title, content, tags1)
	if err != nil {
		fmt.Println("DbInsertpost Error")
	}

	// Send user to loading page
	err = tmpl.ExecuteTemplate(w, "createpostauth", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func feedbackAuthHandler(w http.ResponseWriter, r *http.Request) {
	userId := getUserByCookie(r)

	// Redirect unregistered users
	if userId == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get information from URL
	postRegex := regexp.MustCompile(`post_id=(\d+)`)
	commentRegex := regexp.MustCompile(`comment_id=(\d+)`)
	likeRegex := regexp.MustCompile(`like`)
	dislikeRegex := regexp.MustCompile(`dislike`)

	var postId int
	var commentId int
	var feedback string
	var postnr int

	if commentRegex.MatchString(r.URL.Path) {
		commentID, err := strconv.Atoi(commentRegex.FindStringSubmatch(r.URL.Path)[1])
		if err != nil {
			http.Error(w, "Bad request - invalid comment ID.", http.StatusBadRequest)
			return
		}
		commentId = commentID
		postID, err := strconv.Atoi(postRegex.FindStringSubmatch(r.URL.Path)[1])
		if err != nil {
			http.Error(w, "Bad request - invalid post ID.", http.StatusBadRequest)
			return
		}
		postnr = postID
		if dislikeRegex.MatchString(r.URL.Path) {
			feedback = "-"
		} else if likeRegex.MatchString(r.URL.Path) {
			feedback = "+"
		}

	} else if postRegex.MatchString(r.URL.Path) {
		postID, err := strconv.Atoi(postRegex.FindStringSubmatch(r.URL.Path)[1])
		if err != nil {
			http.Error(w, "Bad request - invalid post ID.", http.StatusBadRequest)
			return
		}
		postId = postID
		postnr = postID
		if dislikeRegex.MatchString(r.URL.Path) {
			feedback = "-"
		} else if likeRegex.MatchString(r.URL.Path) {
			feedback = "+"
		}
	}

	// Add feedback to db
	err1 := database.DbInsertFeedback(postId, commentId, userId, feedback)
	if err1 != nil {
		fmt.Println("DbInsertfeedback Error")
	}
	// Send user to post page
	http.Redirect(w, r, "/post/"+strconv.Itoa(postnr), http.StatusSeeOther)
}
