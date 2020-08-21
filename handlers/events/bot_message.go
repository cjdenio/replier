package events

import (
	"fmt"

	"github.com/cjdenio/replier/db"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleBotDM is called when a user DMs the bot
func HandleBotDM(outer *slackevents.EventsAPICallbackEvent, inner *slackevents.MessageEvent) {
	installation, _ := db.GetInstallation(outer.TeamID)
	client := slack.New(installation.Token)

	_, _, err := client.PostMessage(inner.Channel, slack.MsgOptionText(fmt.Sprintf("you said *%s*", inner.Text), false))
	if err != nil {
		fmt.Println(err)
	}
}
