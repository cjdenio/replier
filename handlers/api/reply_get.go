package api

import (
	"encoding/json"
	"net/http"

	"github.com/cjdenio/replier/db"
)

func HandleAPIReplyGet(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserByToken(r.URL.Query().Get("token"))
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"ok\":false}"))
		return
	}

	encoded, err := json.Marshal(user.Reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"ok\":false}"))
		return
	}

	w.Write(encoded)
}
