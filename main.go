package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"01.kood.tech/git/kretesaak/forum/internal/database"
	"01.kood.tech/git/kretesaak/forum/internal/server"
)

func main() {

	database.InitDb()
	/*
		// Setting up database
		DB := database.DbOpen()

		if DB == nil {
			fmt.Println("500")
		}
		defer DB.Close()
	*/
	defer database.Db.Close()

	fmt.Println("Database open")

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

	//database.TestDbStuff()
	server.StartServer()

}
