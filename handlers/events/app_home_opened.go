package events

import (
	"fmt"
	"log"
	"os"

	"github.com/cjdenio/replier/db"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleAppHomeOpened is fired when the user opens the App Home.
func HandleAppHomeOpened(outer *slackevents.EventsAPICallbackEvent, inner *slackevents.AppHomeOpenedEvent) {
	client := slack.New(os.Getenv("SLACK_TOKEN"))

	user, err := db.GetUser(inner.User)

	var blocks []slack.Block
	if err != nil {
		blocks = []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Hi there! :wave: Please <%s|log in real quick> to get started!", os.Getenv("HOST")+"/login"), false, false),
				nil,
				nil,
			),
		}
	} else {
		blocks = []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", "*Your autoreply message:*", false, false),
				nil,
				nil,
			),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", user.Reply.Message, false, false),
				nil,
				slack.NewAccessory(slack.NewButtonBlockElement("edit_message", "", slack.NewTextBlockObject("plain_text", ":pencil: Edit", true, false))),
			),
			slack.NewDividerBlock(),
		}
	}

	_, err = client.PublishView(inner.User, slack.HomeTabViewRequest{
		Type: "home",
		Blocks: slack.Blocks{
			BlockSet: blocks,
		},
	}, "")

	if err != nil {
		log.Fatal(err)
	}
}
