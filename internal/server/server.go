/*
Server-specific settings
*/
package server

import (
	"html/template"
	"log"
	"net/http"
)
var (
	Port = "8080" // default port
	tmpl *template.Template
)

// Start the server
func StartServer() {

	tmpl = template.Must(template.ParseGlob("templates/*.html"))


	// Serving up the result with mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", registerHandler) // starting endpoint, right now register page
	// Artist endpoint creation

	// Serving up files
	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles/")))) // css serving
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))             // js serving
	err := http.ListenAndServe(":"+Port, mux)
	if err != nil {
		log.Fatalf("Internal server error - %v", http.StatusInternalServerError)
	}

}