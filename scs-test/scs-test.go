package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager

func main() {
	// Initialize a new session manager and configure the session lifetime.
	sessionManager = scs.New()
	// sessionManager := scs.SessionManager{
	// 	contextKey: "Z9JV8hf4AN23Km7u",
	// }
	fmt.Println(sessionManager)
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()
	mux.HandleFunc("/put", putHandler)
	mux.HandleFunc("/get", getHandler)

	addr := ":4000"
	// Wrap your handlers with the LoadAndSave() middleware.
	log.Println("Listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, sessionManager.LoadAndSave(mux)))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	// Store a new key and value in the session data.
	sessionManager.Put(r.Context(), "message", "Hello from a session!")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	// Use the GetString helper to retrieve the string value associated with a
	// key. The zero value is returned if the key does not exist.
	msg := sessionManager.GetString(r.Context(), "message")
	io.WriteString(w, msg)
}
