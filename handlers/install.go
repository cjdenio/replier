package handlers

import (
	"fmt"
	"net/http"
	"os"
)

// HandleInstall redirects the user to the Slack installation
func HandleInstall(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://slack.com/oauth/v2/authorize?scope=im:history,chat:write&client_id=%s&redirect_uri=%s&team=TEHRV8VC", os.Getenv("SLACK_CLIENT_ID"), os.Getenv("HOST")+"/code"), 307)
}
