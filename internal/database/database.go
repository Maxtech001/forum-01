package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const dbfile = "./db/forum.db"

var Db *sql.DB

func InitDb() {
	// Setting up database
	Db = DbOpen()
}

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
func DbInsertUser(user User) error {
	stmt, err := Db.Prepare("INSERT INTO user(id, name, email, password) values(?, ?, ?, ?)")
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
func DbGetUserByIdOrEmail(input string) []User {
	var result []User
	rows, err := Db.Query("SELECT * FROM user WHERE id=? OR email=?", input, input)
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
func DbAuthenticateUser(input, pwd string) bool {
	result := false
	var user User

	err := Db.QueryRow("SELECT password FROM user WHERE id=?", input).Scan(&user.Password) //todo, use count(*) instead
	if err != nil {
		return result
	}
	return CheckPasswordHash(pwd, user.Password)
}
