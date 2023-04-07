/*
Server-specific settings
*/
package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	Port = "8080" // default port
	tmpl *template.Template
	mux  *http.ServeMux
)

// Start the server
func StartServer() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println("Templates:")
	for _, t := range tmpl.Templates() {
		fmt.Println("- " + t.Name())
	}

	// Serving up the result with mux
	mux = http.NewServeMux()
	mux.HandleFunc("/register", registerHandler)         // registration page
	mux.HandleFunc("/registerauth", registerAuthHandler) // registration authentication page
	mux.HandleFunc("/login", loginHandler)               // logging page
	mux.HandleFunc("/loginauth", loginAuthHandler)       // logging authentication page
	mux.HandleFunc("/createpost", createPostHandler)     // creating a post page
	mux.HandleFunc("/logout", logoutHandler)             // logout handler
	mux.HandleFunc("/", mainPageHandler)                 // main page handler

	// Post handling
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		postHandler(w, r)
	})

	// Serving up files
	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles/")))) // css serving
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))             // js serving
	err := http.ListenAndServe(":"+Port, mux)
	if err != nil {
		log.Fatalf("Internal server error - %v", http.StatusInternalServerError)
	}
}
