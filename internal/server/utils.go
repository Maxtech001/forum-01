package server

import (
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
func getMainPageContent(u string) database.Mainpage {
	return database.Mainpage{
		User_id: u,
		Posts:   database.DbGetPosts(),
		Tags:    database.DbGetTags(),
	}
}
