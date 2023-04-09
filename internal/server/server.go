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
	mux  *http.ServeMux
)

// Start the server
func StartServer() {
	// tmpl = template.Must(template.ParseGlob("templates/*.html"))
	// Templating with custom time-helper function
	tmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"formatTime": formatTime,
	}).ParseGlob("templates/*.html"))

	// Serving up the result with mux
	mux = http.NewServeMux()
	mux.HandleFunc("/register", registerHandler)             // registration page
	mux.HandleFunc("/registerauth", registerAuthHandler)     // registration authentication page
	mux.HandleFunc("/login", loginHandler)                   // logging page
	mux.HandleFunc("/loginauth", loginAuthHandler)           // logging authentication page
	mux.HandleFunc("/createpost", createPostHandler)         // creating a post page
	mux.HandleFunc("/logout", logoutHandler)                 // logout handler
	mux.HandleFunc("/", mainPageHandler)                     // main page handler
	mux.HandleFunc("/post/", postHandler)                    // post handling
	mux.HandleFunc("/comment/", commentHandler)              // comment handling
	mux.HandleFunc("/commentauth/", commentAuthHandler)      // comment authentication handling
	mux.HandleFunc("/feedbackauth/", feedbackAuthHandler)    // feedback authentication handling
	mux.HandleFunc("/createpostauth", createPostAuthHandler) //successful post authentication page

	// Serving up files
	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles/")))) // css serving
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))             // js serving
	err := http.ListenAndServe(":"+Port, mux)
	if err != nil {
		log.Fatalf("Internal server error - %v", http.StatusInternalServerError)
	}
}
