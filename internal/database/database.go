package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const dbfile = "./db/forum.db"

var db *sql.DB

// open database
func DbOpen() *sql.DB {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

// insert new user
func dbInsertUser(user User) error {
	stmt, err := db.Prepare("INSERT INTO user(id, name, email, password) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	pwd := HashPassword(user.Password)

	_, err = stmt.Exec(user.Id, user.Name, user.Email, pwd)
	if err != nil {
		return err

	}
	return nil
}

// check if user exists
func dbGetUserByIdOrEmail(input string) []User {
	var result []User
	rows, err := db.Query("SELECT * FROM user WHERE id=? OR email=?", input, input)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			fmt.Println(err)
			return result
		}
		result = append(result, user)
	}
	return result
}

// authenticate by username and password
func dbAuthenticateUser(input, pwd string) bool {
	result := false
	var user User

	err := db.QueryRow("SELECT password FROM user WHERE id=?", input).Scan(&user.Password) //todo, use count(*) instead
	if err != nil {
		return result
	}
	return CheckPasswordHash(pwd, user.Password)
}
