package events

import (
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
				slack.NewTextBlockObject("mrkdwn", "Please log in.", false, false),
				nil,
				nil,
			),
		}
	} else {
		blocks = []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", user.Reply.Message, false, false),
				nil,
				nil,
			),
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
