package events

import (
	"github.com/cjdenio/replier/db"
	"github.com/cjdenio/replier/util"

	"sync"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleMessage handles DMs
func HandleMessage(outer *slackevents.EventsAPICallbackEvent, inner *slackevents.MessageEvent) {
	authedUsers := outer.AuthedUsers
	wg := sync.WaitGroup{}

	wg.Add(len(authedUsers))

	for _, v := range authedUsers {
		go func(userID string) {
			defer wg.Done()

			if userID == inner.User {
				return
			}
			user, err := db.GetUser(userID)

			if err == nil && user.ReplyShouldSend() && !util.IsInArray(user.Reply.Whitelist, inner.User) {
				client := slack.New(user.Token)
				client.PostMessage(inner.Channel, slack.MsgOptionBlocks(
					slack.NewSectionBlock(
						slack.NewTextBlockObject("mrkdwn", user.Reply.Message, false, false),
						nil,
						nil,
					),
					slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", "This is an automatic reply", false, false)),
				))
			}
		}(v)
	}

	wg.Wait()
}
