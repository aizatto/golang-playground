package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"

	"auth"
	"session"

	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

func main() {
	godotenv.Load()
	gob.Register(map[string]interface{}{})
	session.SetupStore()
	mux := http.NewServeMux()
	// mux.HandleFunc("/put", putHandler)
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/auth0/callback", auth.Auth0Callback)

	mux.Handle("/", negroni.New(
		negroni.HandlerFunc(auth.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(getHandler)),
	))
	// n := negroni.Classic()
	// n.Use(
	// 	negroni.HandlerFunc(auth.IsAuthenticated))
	// // n.UseHandler(mux)

	addr := ":3000"
	log.Println("Listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	// Store a new key and value in the session data.
	session, err := session.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["foo"] = "bar"
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	// Use the GetString helper to retrieve the string value associated with a
	// key. The zero value is returned if the key does not exist.
	// session, err := session.Store.Get(r, "session-name")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// val := session.Values["foo"].(string)
	line := fmt.Sprintf("Req: %s %s\n", r.Host, r.URL.Path)
	io.WriteString(w, line)
}
