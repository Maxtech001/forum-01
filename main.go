package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"01.kood.tech/git/kretesaak/forum/internal/database"
	"01.kood.tech/git/kretesaak/forum/internal/server"
)

func main() {
	db, err := database.InitDB()
	if err != nil {

		log.Fatalf("Internal server error - %v", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	fmt.Println("Database connection open")

	go func() {
		for {
			database.DbDeleteExpiredCookies()
			time.Sleep(1 * time.Hour)
		}
	}()

	// Starting server
	if len(os.Args) == 1 {
	} else if len(os.Args) == 2 {
		server.Port = os.Args[1]
	} else {
		log.Fatalf("Internal server error - %v", http.StatusInternalServerError)
		return
	}
	fmt.Println("Server starting...")
	fmt.Println("\nOpen http://localhost:" + " with the specified port in browser")
	fmt.Println("\nCtrl + C to close server")

	server.StartServer()
}
