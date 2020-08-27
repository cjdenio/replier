package handlers

import (
	"log"
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
		_, err = w.Write([]byte("Something went wrong; please try again. :("))
		if err != nil {
			log.Println(err)
		}
		return
	}
	if resp.AuthedUser.AccessToken != "" {
		err = db.AddUser(db.User{
			Token:  resp.AuthedUser.AccessToken,
			UserID: resp.AuthedUser.ID,
			Scopes: strings.Split(resp.AuthedUser.Scope, ","),
			TeamID: resp.Team.ID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("Something went wrong on our end. Please try again in a little bit."))
			if err != nil {
				log.Println(err)
			}
		}
	}

	if resp.AccessToken != "" {
		err := db.AddInstallation(db.Installation{
			Token:  resp.AccessToken,
			Scopes: strings.Split(resp.Scope, ","),
			TeamID: resp.Team.ID,
			BotID:  resp.BotUserID,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("Something went wrong on our end. Please try again in a little bit."))
			if err != nil {
				log.Println(err)
			}
		}
	}
	w.Header().Add("Content-Type", "text/html")
	_, err = w.Write([]byte("<h1 style='font-family:sans-serif'>You're logged in!</h1><p style='font-family:sans-serif'>You can now head on back to Slack.</p>"))
	if err != nil {
		log.Println(err)
	}

	err = util.UpdateAppHome(resp.AuthedUser.ID, resp.Team.ID)
	if err != nil {
		log.Println(err)
	}

	err = util.SendWelcomeMessage(resp.Team.ID, resp.AuthedUser.ID)
	if err != nil {
		log.Println(err)
	}
}
