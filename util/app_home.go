package util

import (
	"fmt"
	"os"

	"github.com/cjdenio/replier/db"
	"github.com/slack-go/slack"
)

// NewInputBlock is an input block that contains the DispatchAction field
type NewInputBlock struct {
	Type           slack.MessageBlockType `json:"type"`
	BlockID        string                 `json:"block_id,omitempty"`
	Label          *slack.TextBlockObject `json:"label"`
	Element        slack.BlockElement     `json:"element"`
	Hint           *slack.TextBlockObject `json:"hint,omitempty"`
	Optional       bool                   `json:"optional,omitempty"`
	DispatchAction bool                   `json:"dispatch_action"`
}

func (s NewInputBlock) BlockType() slack.MessageBlockType {
	return s.Type
}

func UpdateAppHome(userID, teamID string) error {
	installation, err := db.GetInstallation(teamID)
	if err != nil {
		fmt.Println(err)
	}
	client := slack.New(installation.Token)

	user, err := db.GetUser(userID)

	replyMode := user.Reply.Mode
	if replyMode == "" {
		replyMode = db.ReplyModeManual
	}

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
					Text:     slack.NewTextBlockObject("plain_text", ":bust_in_silhouette: Log in", true, false),
					ActionID: "login",
					URL:      os.Getenv("HOST") + "/login",
				}),
			),
		}
	} else {
		replyActive := user.ReplyShouldSend()

		replyActiveText := ":x: Your autoreply is *off*."
		if replyActive {
			replyActiveText = ":heavy_check_mark: Your autoreply is *on*!"
		}
		if user.Reply.Message == "" {
			replyActiveText = ":x: Your autoreply is *off* because you haven't set a message."
		}

		var replyActiveAccessory *slack.Accessory

		if user.Reply.Mode == db.ReplyModeManual {
			replyActiveAccessory = slack.NewAccessory(&slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     slack.NewTextBlockObject("plain_text", map[bool]string{true: "Turn off", false: "Turn on"}[user.Reply.Active], false, false),
				ActionID: "reply_toggle",
			})
		}

		buttonStyles := map[string]slack.Style{
			"manual":   "",
			"date":     "",
			"presence": "",
		}
		if replyMode == db.ReplyModeManual {
			buttonStyles["manual"] = slack.StylePrimary
		} else if replyMode == db.ReplyModeDate {
			buttonStyles["date"] = slack.StylePrimary
		} else if replyMode == db.ReplyModePresence {
			buttonStyles["presence"] = slack.StylePrimary
		}

		blocks = []slack.Block{
			&NewInputBlock{
				Type:           slack.MBTInput,
				Label:          slack.NewTextBlockObject("plain_text", "Your autoreply message", false, false),
				BlockID:        "message",
				DispatchAction: true,
				Element: &slack.PlainTextInputBlockElement{
					Type:         slack.METPlainTextInput,
					Multiline:    true,
					ActionID:     "message",
					InitialValue: user.Reply.Message,
				},
			},
			slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", ":sparkles: *Fun fact:* if you put `@person` in the message, it'll get replaced by the actual message sender's name!", false, false)),
			slack.NewActionBlock("", slack.NewButtonBlockElement("edit_message", "", slack.NewTextBlockObject("plain_text", ":gear: Settings", true, false))),
			slack.NewDividerBlock(),
			slack.NewActionBlock("", slack.NewButtonBlockElement("mode-manual", "", slack.NewTextBlockObject("plain_text", "Manual", false, false)).WithStyle(buttonStyles["manual"]), slack.NewButtonBlockElement("mode-date", "", slack.NewTextBlockObject("plain_text", "Date Range", false, false)).WithStyle(buttonStyles["date"]), slack.NewButtonBlockElement("mode-presence", "", slack.NewTextBlockObject("plain_text", "Presence", false, false)).WithStyle(buttonStyles["presence"])),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", replyActiveText, false, false),
				nil,
				replyActiveAccessory,
			),
			/*slack.NewDividerBlock(),
			slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", "Replier is open-source on <https://github.com/cjdenio/replier|GitHub>!", false, false)),*/
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
