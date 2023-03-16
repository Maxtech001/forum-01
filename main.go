package main

import (
	"fmt"
	"os"
	"log"
	"net/http"

	"01.kood.tech/git/kretesaak/forum/internal/database"
	"01.kood.tech/git/kretesaak/forum/internal/server"
)


func main() {

	// Setting up database
	db := database.DbOpen()

	if db == nil {
		fmt.Println("500")
	}
	defer db.Close()

	// Starting server
	if len(os.Args) == 1 {
		} else if len(os.Args) == 2 {
			server.Port = os.Args[1]
		} else {
			log.Fatalf("Internal server error - %v", http.StatusInternalServerError)
			return
		}
	fmt.Println("Server starting on port", server.Port)
	fmt.Println("\nOpen http://localhost:" + server.Port + "/ in browser")
	fmt.Println("\nCtrl + C to close server")

	
	server.StartServer()

	// register new user
	/*
		var newuser user
		newuser.Id = "margusaid"
		newuser.Name = "Margus Aid"
		newuser.Email = "margus.axx@gmail.com"
		newuser.Password = "12345"

		err := dbInsertUser(newuser)

		if err != nil {
			fmt.Println("500", err)
		}
	*/
	
	/*
	// get user data by username or e-mail address, returns user struct
	//fmt.Println(dbGetUserByIdOrEmail("n00bh4ck3r"))
	fmt.Println(dbGetUserByIdOrEmail("margus.axx@gmail.com"))

	// authenticate user by email/username and password, returns boolean
	//fmt.Println(dbAuthenticateUser("n00bh4ck3r", "qwerty"))
	fmt.Println(dbAuthenticateUser("margusaid", "12345"))
	*/

}
