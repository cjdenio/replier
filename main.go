package main

import (
	"net/http"
	"os"

	"github.com/cjdenio/replier/db"
	"github.com/cjdenio/replier/handlers"
)

func main() {
	db.Connect()

	http.HandleFunc("/slack/events", handlers.HandleEvents)
	http.HandleFunc("/slack/interactivity", handlers.HandleInteractivity)
	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/code", handlers.HandleOAuthCode)

	http.ListenAndServe(":"+os.Getenv("PORT"), http.DefaultServeMux)
}
