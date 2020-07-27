package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/cjdenio/replier/db"
	"github.com/slack-go/slack"
	//"github.com/cjdenio/replier/db"
)

// HandleOAuthCode handles the OAuth redirect
func HandleOAuthCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	resp, err := slack.GetOAuthV2Response(&http.Client{}, os.Getenv("SLACK_CLIENT_ID"), os.Getenv("SLACK_CLIENT_SECRET"), code, os.Getenv("HOST")+"/code")
	if err != nil {
		log.Fatal(err)
	}
	db.AddUser(db.User{
		Token:  resp.AuthedUser.AccessToken,
		UserID: resp.AuthedUser.ID,
	})
}
