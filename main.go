package main

import (
	"net/http"

	"github.com/cjdenio/replier/db"
	"github.com/cjdenio/replier/handlers"
)

func main() {
	db.Connect()

	http.HandleFunc("/slack/events", handlers.HandleEvents)
	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/code", handlers.HandleOAuthCode)

	http.ListenAndServe(":3000", http.DefaultServeMux)
}
