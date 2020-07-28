package events

import (
	"fmt"

	"github.com/cjdenio/replier/util"
	"github.com/slack-go/slack/slackevents"
)

// HandleAppHomeOpened is fired when the user opens the App Home.
func HandleAppHomeOpened(outer *slackevents.EventsAPICallbackEvent, inner *slackevents.AppHomeOpenedEvent) {
	if err := util.UpdateAppHome(inner.User); err != nil {
		fmt.Println(err)
	}
}
