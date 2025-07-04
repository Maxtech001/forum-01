package server

import (
	"errors"
	"net/http"
	"time"

	"01.kood.tech/git/kretesaak/forum/internal/database"
	_ "github.com/gofrs/uuid"
)

// Reformatting time helper in front-end
func formatTime(t string) string {
	tc, _ := time.Parse(time.RFC3339, t)
	return tc.Format("2 January 2006")
}

// Returning main page content
func getMainPageContent(u string, qparams map[string][]string) database.Mainpage {
	return database.Mainpage{
		User_id: u,
		Posts:   database.DbGetPosts(u, qparams),
		Tags:    database.DbGetTags(),
	}
}

// Returning create post page content
func getCreatePostPageContent(u string) database.Createpost {
	return database.Createpost{
		User_id: u,
		Tags:    database.DbGetTags(),
	}
}

// Returning post page content
func getPostPageContent(pID int, u string) (error, database.Postpage) {
	post := database.DbGetSinglePost(pID, u)

	if post.Id == 0 {
		return errors.New("No such post"), database.Postpage{}
	}
	return nil, database.Postpage{
		User_id: u,
		Post:    post,
	}
}

// Returning a new created user from database
func newUser(u, e, p string) database.User {
	return database.User{
		Id:       u,
		Email:    e,
		Password: p,
	}
}

func getUserByCookie(r *http.Request) string {
	result := ""
	cookie, err := r.Cookie("session")
	if err != nil {
		return result
	}
	result = database.DbGetUserByCookie(cookie.Value)
	return result
}
