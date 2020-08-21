package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/cjdenio/replier/db"
	"github.com/cjdenio/replier/util"
	"github.com/slack-go/slack"
	//"github.com/cjdenio/replier/db"
)

// HandleOAuthCode handles the OAuth redirect
func HandleOAuthCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	resp, err := slack.GetOAuthV2Response(&http.Client{}, os.Getenv("SLACK_CLIENT_ID"), os.Getenv("SLACK_CLIENT_SECRET"), code, os.Getenv("HOST")+"/code")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Something went wrong. :("))
		return
	}
	if resp.AuthedUser.AccessToken != "" {
		db.AddUser(db.User{
			Token:  resp.AuthedUser.AccessToken,
			UserID: resp.AuthedUser.ID,
			Scopes: strings.Split(resp.AuthedUser.Scope, ","),
			TeamID: resp.Team.ID,
		})
	}

	if resp.AccessToken != "" {
		db.AddInstallation(db.Installation{
			Token:  resp.AccessToken,
			Scopes: strings.Split(resp.Scope, ","),
			TeamID: resp.Team.ID,
			BotID:  resp.BotUserID,
		})
	}
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<h1 style='font-family:sans-serif'>You're logged in!</h1><p style='font-family:sans-serif'>You can now head on back to Slack.</p>"))

	util.UpdateAppHome(resp.AuthedUser.ID, resp.Team.ID)

	util.SendWelcomeMessage(resp.Team.ID, resp.AuthedUser.ID)
}
