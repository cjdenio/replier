package events

import (
	"fmt"
	"strings"

	"github.com/cjdenio/replier/db"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleBotDM is called when a user DMs the bot
func HandleBotDM(outer *slackevents.EventsAPICallbackEvent, inner *slackevents.MessageEvent) {
	installation, _ := db.GetInstallation(outer.TeamID)
	client := slack.New(installation.Token)

	token, _ := db.GetUserAPIToken(inner.User)

	if strings.Contains(strings.ToLower(inner.Text), "api") {
		_, _, err := client.PostMessage(inner.Channel, slack.MsgOptionText(fmt.Sprintf("Your API token is `%s`! Guard it well, as it grants access to your Replier settings.\n\n<https://github.com/cjdenio/replier/wiki|API docs>", token), false))
		if err != nil {
			fmt.Println(err)
		}
	}
}
