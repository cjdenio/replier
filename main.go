package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cjdenio/replier/db"
	"github.com/cjdenio/replier/handlers"
)

func main() {
	db.Connect()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://github.com/cjdenio/replier", http.StatusMovedPermanently)
	})
	http.HandleFunc("/slack/events", handlers.HandleEvents)
	http.HandleFunc("/slack/interactivity", handlers.HandleInteractivity)
	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/install", handlers.HandleInstall)
	http.HandleFunc("/code", handlers.HandleOAuthCode)

	err := http.ListenAndServe(":"+os.Getenv("PORT"), http.DefaultServeMux)

	log.Fatal(err)
}
