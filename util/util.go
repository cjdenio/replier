package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/cjdenio/replier/db"
	"github.com/slack-go/slack"
)

// UpdateAppHome updates the App Home for the given user
func UpdateAppHome(userID string) error {
	client := slack.New(os.Getenv("SLACK_TOKEN"))

	user, err := db.GetUser(userID)

	var blocks []slack.Block
	if err != nil {
		fmt.Println(err)
		blocks = []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Hi there! :wave: Please <%s|log in real quick> to get started!", os.Getenv("HOST")+"/login"), false, false),
				nil,
				nil,
			),
		}
	} else {
		replyMessage := user.Reply.Message

		if replyMessage == "" {
			replyMessage = "*You haven't set up an autoreply yet.* Click that button over there :arrow_right: to get started!"
		}

		replyActiveText := ":x: Your autoreply isn't active. That means that people will *not* receive it when they DM you."
		replyToggleButtonText := "Turn On"

		if user.Reply.Active {
			replyActiveText = ":heavy_check_mark: Your autoreply is active! That means that people *will* receive it when they attempt to DM you."
			replyToggleButtonText = "Turn Off"
		}

		blocks = []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", "*Your autoreply message:*", false, false),
				nil,
				nil,
			),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", replyMessage, false, false),
				nil,
				slack.NewAccessory(slack.NewButtonBlockElement("edit_message", "", slack.NewTextBlockObject("plain_text", ":pencil: Edit", true, false))),
			),
			slack.NewDividerBlock(),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", replyActiveText, false, false),
				nil,
				slack.NewAccessory(slack.NewButtonBlockElement(
					"reply_toggle",
					"",
					slack.NewTextBlockObject("plain_text", replyToggleButtonText, false, false),
				)),
			),
		}
	}

	_, err = client.PublishView(userID, slack.HomeTabViewRequest{
		Type: "home",
		Blocks: slack.Blocks{
			BlockSet: blocks,
		},
	}, "")

	if err != nil {
		return err
	}

	return nil
}

// VerifySlackRequest verifies a Slack request
func VerifySlackRequest(r *http.Request, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(os.Getenv("SLACK_SIGNING_SECRET")))

	body = append([]byte(r.Header.Get("X-Slack-Request-Timestamp")+":"), body...)
	body = append([]byte("v0:"), body...)

	mac.Write(body)

	return hmac.Equal([]byte("v0="+hex.EncodeToString(mac.Sum(nil))), []byte(r.Header.Get("X-Slack-Signature")))
}
