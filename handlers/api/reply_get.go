package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cjdenio/replier/db"
)

func HandleAPIReplyGet(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserByToken(r.URL.Query().Get("token"))
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte("{\"ok\":false}"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	encoded, err := json.Marshal(user.Reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("{\"ok\":false}"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	_, err = w.Write(encoded)
	if err != nil {
		log.Println(err)
	}
}
