package api

import "net/http"

func HandleAPIReplyGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
