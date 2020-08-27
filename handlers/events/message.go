package events

import (
	"fmt"
	"log"
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

			if userID == inner.User || inner.BotID != "" || inner.User == "USLACKBOT" {
				return
			}
			user, err := db.GetUser(userID)
			lastPostedOn := db.GetConversationLastPostedOn(inner.Channel, userID)

			if err == nil && user.ReplyShouldSend() && !util.IsInArray(user.Reply.Whitelist, inner.User) && time.Since(lastPostedOn).Minutes() > 15 {
				client := slack.New(user.Token)
				_, _, err = client.PostMessage(inner.Channel, slack.MsgOptionBlocks(
					slack.NewSectionBlock(
						slack.NewTextBlockObject("mrkdwn", util.TransformUserReply(user.Reply.Message, inner.User), false, false),
						nil,
						nil,
					),
					slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<slack://app?team=%s&id=%s&tab=home|This is an automatic reply>", outer.TeamID, outer.APIAppID), false, false)),
				), slack.MsgOptionText(util.TransformUserReply(user.Reply.Message, inner.User), false))
				if err != nil {
					log.Println(err)
				}
				if err = db.SetConversationLastPostedOn(inner.Channel, userID, time.Now()); err != nil {
					log.Println(err)
				}
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

			if userID == inner.User || inner.BotID != "" {
				return
			}
			user, err := db.GetUser(userID)

			timestampToReplyTo := inner.ThreadTimeStamp

			if timestampToReplyTo == "" {
				timestampToReplyTo = inner.TimeStamp
			}

			if err == nil && user.ReplyShouldSend() && strings.Contains(inner.Text, fmt.Sprintf("<@%s>", userID)) && !util.IsInArray(user.Reply.Whitelist, inner.User) {
				client := slack.New(user.Token)
				_, _, err = client.PostMessage(inner.Channel, slack.MsgOptionBlocks(
					slack.NewSectionBlock(
						slack.NewTextBlockObject("mrkdwn", util.TransformUserReply(user.Reply.Message, inner.User), false, false),
						nil,
						nil,
					),
					slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<slack://app?team=%s&id=%s&tab=home|This is an automatic reply>", outer.TeamID, outer.APIAppID), false, false)),
				), slack.MsgOptionText(util.TransformUserReply(user.Reply.Message, inner.User), false), slack.MsgOptionTS(timestampToReplyTo))
				if err != nil {
					log.Println(err)
				}
			}
		}(v)
	}

	wg.Wait()
}
