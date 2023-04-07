package registration

import (
	"01.kood.tech/git/kretesaak/forum/internal/database"
)

// Init new user
func NewUser(u, e, p string) database.User {
	return database.User{
		Id:       u,
		Email:    e,
		Password: p,
	}
}
