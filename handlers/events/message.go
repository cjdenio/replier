package events

import (
	"github.com/cjdenio/replier/db"
	//"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleMessage handles DMs
func HandleMessage(outer *slackevents.EventsAPICallbackEvent, inner *slackevents.MessageEvent) {
	authedUsers := outer.AuthedUsers
	for _, v := range authedUsers {
		go func(userID string) {
			user := db.GetUser(userID)

			if user.HasActiveReply() {
				//client := slack.New(user.Token)
			}
		}(v)
	}
}
