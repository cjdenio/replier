package events

import (
	"fmt"
	"strings"
	"time"

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
			lastPostedOn := db.GetConversationLastPostedOn(inner.Channel, userID)

			if err == nil && user.ReplyShouldSend() && !util.IsInArray(user.Reply.Whitelist, inner.User) && time.Now().Sub(lastPostedOn).Minutes() > 15 {
				client := slack.New(user.Token)
				client.PostMessage(inner.Channel, slack.MsgOptionBlocks(
					slack.NewSectionBlock(
						slack.NewTextBlockObject("mrkdwn", user.Reply.Message, false, false),
						nil,
						nil,
					),
					slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<slack://app?team=%s&id=%s&tab=home|This is an automatic reply>", outer.TeamID, outer.APIAppID), false, false)),
				))
				db.SetConversationLastPostedOn(inner.Channel, userID, time.Now())
			}
		}(v)
	}

	wg.Wait()
}

// HandleMessageNonDM handles non-DM messages
func HandleMessageNonDM(outer *slackevents.EventsAPICallbackEvent, inner *slackevents.MessageEvent) {
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

			timestampToReplyTo := inner.ThreadTimeStamp

			if timestampToReplyTo == "" {
				timestampToReplyTo = inner.TimeStamp
			}

			if err == nil && user.ReplyShouldSend() && strings.Contains(inner.Text, fmt.Sprintf("<@%s>", userID)) {
				client := slack.New(user.Token)
				client.PostMessage(inner.Channel, slack.MsgOptionBlocks(
					slack.NewSectionBlock(
						slack.NewTextBlockObject("mrkdwn", user.Reply.Message, false, false),
						nil,
						nil,
					),
					slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<slack://app?team=%s&id=%s&tab=home|This is an automatic reply>", outer.TeamID, outer.APIAppID), false, false)),
				), slack.MsgOptionTS(timestampToReplyTo))
			}
		}(v)
	}

	wg.Wait()
}
