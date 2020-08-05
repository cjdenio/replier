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

type HeaderBlock struct {
	Type string                 `json:"type"`
	Text *slack.TextBlockObject `json:"text"`
}

func (b HeaderBlock) BlockType() slack.MessageBlockType {
	return slack.MessageBlockType(b.Type)
}

// UpdateAppHome updates the App Home for the given user
func UpdateAppHome(userID string) error {
	client := slack.New(os.Getenv("SLACK_TOKEN"))

	user, err := db.GetUser(userID)

	needsToLogin := false

	if err != nil || user.Token == "" {
		needsToLogin = true
	} else if _, err := slack.New(user.Token).AuthTest(); err != nil {
		needsToLogin = true
	}

	var blocks []slack.Block
	if needsToLogin {
		blocks = []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Hi there! :wave: Please <%s|log in real quick> to get started!", os.Getenv("HOST")+"/login"), false, false),
				nil,
				slack.NewAccessory(&slack.ButtonBlockElement{
					Type:     slack.METButton,
					Text:     slack.NewTextBlockObject("plain_text", ":bust_in_silhouette: Login", true, false),
					ActionID: "login",
					URL:      os.Getenv("HOST") + "/login",
				}),
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
				slack.NewAccessory(slack.NewButtonBlockElement("edit_message", "", slack.NewTextBlockObject("plain_text", ":pencil: Edit settings", true, false))),
			),
			slack.NewDividerBlock(),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", replyActiveText, false, false),
				nil,
				slack.NewAccessory(&slack.ButtonBlockElement{
					Type:     slack.METButton,
					Text:     slack.NewTextBlockObject("plain_text", replyToggleButtonText, false, false),
					ActionID: "reply_toggle",
				}),
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

// IsInArray checks if the value is in the array
func IsInArray(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func GetUserTimezone(userID string) (string, error) {
	user, err := db.GetUser(userID)
	if err != nil {
		return "", err
	}

	client := slack.New(user.Token)
	slackUser, err := client.GetUserInfo(user.UserID)
	if err != nil {
		return "", err
	}

	return slackUser.TZ, nil
}
