package handlers

import (
	"fmt"
	"net/http"
	"os"
)

// HandleLogin redirects the user to the Slack login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://slack.com/oauth/v2/authorize?user_scope=im:history,mpim:history,chat:write&client_id=%s&redirect_uri=%s", os.Getenv("HOST")+"/code", os.Getenv("SLACK_REDIRECT_URI")), 307)
}
