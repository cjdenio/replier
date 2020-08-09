package handlers

import (
	"fmt"
	"net/http"
	"os"
)

// HandleLogin redirects the user to the Slack login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://slack.com/oauth/v2/authorize?scope=im:history,chat:write&user_scope=im:history,mpim:history,channels:history,groups:history,chat:write,users:read&client_id=%s&redirect_uri=%s", os.Getenv("SLACK_CLIENT_ID"), os.Getenv("HOST")+"/code"), 307)
}
