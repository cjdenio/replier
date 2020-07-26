package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://slack.com/oauth/v2/authorize?user_scope=im:history&client_id=%s&redirect_uri=%s", os.Getenv("SLACK_CLIENT_ID"), os.Getenv("SLACK_REDIRECT_URI")), 307)
}
