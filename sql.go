package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const dbfile = "./db/forum.db"

var db *sql.DB

// open database
func dbOpen() *sql.DB {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

// insert new user
func dbInsertUser(user user) error {
	stmt, err := db.Prepare("INSERT INTO user(id, name, email, password) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		return err

	}
	return nil
}

// check if user exists
func dbGetUserByIdOrEmail(input string) []user {
	var result []user
	rows, err := db.Query("SELECT * FROM user WHERE id=? OR email=?", input, input)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var user user
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			fmt.Println(err)
			return result
		}
		result = append(result, user)
	}
	return result
}

// authenticate by email or username and password
func dbAuthenticateUser(input, pwd string) bool {
	result := false
	var user user
	err := db.QueryRow("SELECT id FROM user WHERE (id=? OR email=?) AND password=?", input, input, pwd).Scan(&user.Id)
	if err != nil {
		return result
	}
	if user.Id != "" {
		result = true
	}
	return result
}
