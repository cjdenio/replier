package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

			client := slack.New(user.Token)
			slackUser, err := client.GetUserInfo(user.UserID)

			var initialStartDate string
			var initialEndDate string

			if user.Reply.Start != (time.Time{}) {
				initialStartDate = user.Reply.Start.Format("2006-01-02")
			}

			if user.Reply.End != (time.Time{}) {
				initialEndDate = user.Reply.End.Format("2006-01-02")
			}

			if err != nil {
				blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("_Psst!_ Wanna set start/end dates for your autoreply? <%s|Click here to add that permission!>", os.Getenv("HOST")+"/login"), false, false), nil, nil))
			} else {
				blocks = append(blocks, &util.HeaderBlock{
					Type: "header",
					Text: slack.NewTextBlockObject("plain_text", ":calendar: Dates", true, false),
				},
					&slack.InputBlock{
						Type:     slack.MBTInput,
						BlockID:  "start",
						Label:    slack.NewTextBlockObject("plain_text", "Start Date", false, false),
						Optional: true,
						Element: &slack.DatePickerBlockElement{
							Type:        slack.METDatepicker,
							ActionID:    "start",
							InitialDate: initialStartDate,
						},
					},
					&slack.InputBlock{
						Type:     slack.MBTInput,
						BlockID:  "end",
						Label:    slack.NewTextBlockObject("plain_text", "End Date", false, false),
						Optional: true,
						Element: &slack.DatePickerBlockElement{
							Type:        slack.METDatepicker,
							ActionID:    "end",
							InitialDate: initialEndDate,
						},
					},
					slack.NewContextBlock(
						"",
						slack.NewTextBlockObject(
							"mrkdwn",
							fmt.Sprintf("These dates will be evaluted in your timezone: *%s*", slackUser.TZ),
							false,
							false,
						),
					),
					slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", ":information_source: You will still need to enable the autoreply for it to be sent.", false, false), nil, nil))
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
				log.Fatal(err)
			}
		case "reply_toggle":
			db.ToggleReplyActive(parsed.User.ID)
			err := util.UpdateAppHome(parsed.User.ID, parsed.Team.ID)
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
