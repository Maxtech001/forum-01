ppackage main

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

	// ✅ Set default port if not provided via command-line
	if len(os.Args) == 2 {
		server.Port = os.Args[1]
	} else {
		server.Port = "8080"
	}

	fmt.Println("Server starting...")
	fmt.Printf("\nOpen http://localhost:%s in your browser\n", server.Port)
	fmt.Println("\nPress Ctrl + C to stop the server")

	server.StartServer()
}
