package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cjdenio/replier/db"
	"github.com/cjdenio/replier/handlers"
	"github.com/cjdenio/replier/handlers/api"
	"github.com/gorilla/mux"
)

func main() {
	db.Connect()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://github.com/cjdenio/replier", http.StatusMovedPermanently)
	})
	r.HandleFunc("/slack/events", handlers.HandleEvents).Methods("POST")
	r.HandleFunc("/slack/interactivity", handlers.HandleInteractivity).Methods("POST")
	r.HandleFunc("/login", handlers.HandleLogin)
	r.HandleFunc("/install", handlers.HandleInstall)
	r.HandleFunc("/code", handlers.HandleOAuthCode)

	r.HandleFunc("/api/reply", api.HandleAPIReplyGet).Methods("GET")

	err := http.ListenAndServe(":"+os.Getenv("PORT"), r)

	log.Fatal(err)
}
