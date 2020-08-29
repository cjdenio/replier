package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/cjdenio/replier/db"
	"github.com/cjdenio/replier/util"
	"github.com/slack-go/slack"
)

// HandleInteractivity handles interactions in Slack
func HandleInteractivity(w http.ResponseWriter, r *http.Request) {
	buf, _ := ioutil.ReadAll(r.Body)
	r.Form, _ = url.ParseQuery(string(buf))

	if !util.VerifySlackRequest(r, buf) {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Not verified :("))
		if err != nil {
			log.Println(err)
		}
		return
	}

	var parsed slack.InteractionCallback
	err := json.Unmarshal([]byte(r.Form.Get("payload")), &parsed)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Invalid JSON payload"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	_, err = w.Write(nil)
	if err != nil {
		log.Println(err)
	}

	if parsed.Type == slack.InteractionTypeBlockActions {
		switch parsed.ActionCallback.BlockActions[0].ActionID {
		case "edit_message":
			user, _ := db.GetUser(parsed.User.ID)

			blocks := []slack.Block{
				&slack.InputBlock{
					Type:    slack.MBTInput,
					BlockID: "message",
					Label:   slack.NewTextBlockObject("plain_text", "Message", false, false),
					Element: &slack.PlainTextInputBlockElement{
						Type:         slack.METPlainTextInput,
						ActionID:     "message",
						Multiline:    true,
						InitialValue: user.Reply.Message,
					},
					Optional: true,
				},
				slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", ":sparkles: *Fun fact:* if you put `@person` in the message, it'll get replaced by the actual message sender's name!", false, false)),
				&slack.InputBlock{
					Type:    slack.MBTInput,
					BlockID: "whitelist",
					Label:   slack.NewTextBlockObject("plain_text", "Whitelist", false, false),
					Element: &slack.MultiSelectBlockElement{
						Type:         "multi_users_select",
						InitialUsers: user.Reply.Whitelist,
						ActionID:     "whitelist",
						Placeholder:  slack.NewTextBlockObject("plain_text", "Select some...", false, false),
					},
					Optional: true,
				},
				slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", "These people will _not_ receive your autoreply in DMs or public channels, even if it's enabled.", false, false)),
			}

			installation, err := db.GetInstallation(parsed.Team.ID)
			if err != nil {
				fmt.Println(err)
			}

			botClient := slack.New(installation.Token)
			_, err = botClient.OpenView(parsed.TriggerID, slack.ModalViewRequest{
				Type:       "modal",
				Title:      slack.NewTextBlockObject("plain_text", "Edit Settings", false, false),
				CallbackID: "edit_message",
				Blocks: slack.Blocks{
					BlockSet: blocks,
				},
				Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
				Submit: slack.NewTextBlockObject("plain_text", "Save", false, false),
			})

			if err != nil {
				log.Println(err)
			}
		case "reply_toggle":
			db.ToggleReplyActive(parsed.User.ID)
			err := util.UpdateAppHome(parsed.User.ID, parsed.Team.ID)
			if err != nil {
				log.Println(err)
			}
		case "mode-manual":
			err = db.SetReplyMode(parsed.User.ID, db.ReplyModeManual)
			if err != nil {
				log.Println(err)
			}
			err = util.UpdateAppHome(parsed.User.ID, parsed.Team.ID)
			if err != nil {
				log.Println(err)
			}
		case "mode-date":
			err = db.SetReplyMode(parsed.User.ID, db.ReplyModeDate)
			if err != nil {
				log.Println(err)
			}
			err = util.UpdateAppHome(parsed.User.ID, parsed.Team.ID)
			if err != nil {
				log.Println(err)
			}
		case "mode-presence":
			err = db.SetReplyMode(parsed.User.ID, db.ReplyModePresence)
			if err != nil {
				log.Println(err)
			}
			err = util.UpdateAppHome(parsed.User.ID, parsed.Team.ID)
			if err != nil {
				log.Println(err)
			}
		}
	} else if parsed.Type == slack.InteractionTypeViewSubmission {
		switch parsed.View.CallbackID {
		case "edit_message":
			desiredMessage := parsed.View.State.Values["message"]["message"].Value
			whitelist := parsed.View.State.Values["whitelist"]["whitelist"].SelectedUsers

			tz, err := util.GetUserTimezone(parsed.User.ID)
			if err != nil {
				fmt.Println(err)
			}

			err = db.SetUserMessage(parsed.User.ID, desiredMessage)
			if err != nil {
				log.Println(err)
			}
			err = db.SetUserWhitelist(parsed.User.ID, whitelist)
			if err != nil {
				log.Println(err)
			}

			loc, _ := time.LoadLocation(tz)

			startDate, _ := time.ParseInLocation("2006-01-02", parsed.View.State.Values["start"]["start"].SelectedDate, loc)
			endDate, _ := time.ParseInLocation("2006-01-02", parsed.View.State.Values["end"]["end"].SelectedDate, loc)

			err = db.SetUserDates(startDate, endDate, parsed.User.ID)
			if err != nil {
				log.Println(err)
			}
			err = util.UpdateAppHome(parsed.User.ID, parsed.Team.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
