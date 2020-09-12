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

	if parsed.Type == slack.InteractionTypeBlockActions {
		_, err = w.Write(nil)
		if err != nil {
			log.Println(err)
		}

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
			installation, _ := db.GetInstallation(parsed.Team.ID)
			user, _ := db.GetUser(parsed.User.ID)
			client := slack.New(installation.Token)

			startDate := user.Reply.Start.Format("2006-01-02")
			if (user.Reply.Start == time.Time{}) {
				startDate = ""
			}
			endDate := user.Reply.End.Format("2006-01-02")
			if (user.Reply.End == time.Time{}) {
				endDate = ""
			}

			_, err = client.OpenView(parsed.TriggerID, slack.ModalViewRequest{
				Type:       "modal",
				Title:      slack.NewTextBlockObject("plain_text", "Date Range", false, false),
				CallbackID: "date_range",
				Blocks: slack.Blocks{
					BlockSet: []slack.Block{
						slack.InputBlock{
							Type:     "input",
							BlockID:  "start",
							Optional: true,
							Label:    slack.NewTextBlockObject("plain_text", "Start Date", false, false),
							Element: slack.DatePickerBlockElement{
								Type:        "datepicker",
								InitialDate: startDate,
								ActionID:    "start",
							},
						},
						slack.InputBlock{
							Type:     "input",
							BlockID:  "end",
							Optional: true,
							Label:    slack.NewTextBlockObject("plain_text", "End Date", false, false),
							Element: slack.DatePickerBlockElement{
								Type:        "datepicker",
								InitialDate: endDate,
								ActionID:    "end",
							},
						},
					},
				},
				Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
				Submit: slack.NewTextBlockObject("plain_text", "Save", false, false),
			})
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
			_, err = w.Write(nil)
			if err != nil {
				log.Println(err)
			}

			desiredMessage := parsed.View.State.Values["message"]["message"].Value
			whitelist := parsed.View.State.Values["whitelist"]["whitelist"].SelectedUsers

			err = db.SetUserMessage(parsed.User.ID, desiredMessage)
			if err != nil {
				log.Println(err)
			}
			err = db.SetUserWhitelist(parsed.User.ID, whitelist)
			if err != nil {
				log.Println(err)
			}

			err = util.UpdateAppHome(parsed.User.ID, parsed.Team.ID)
			if err != nil {
				log.Println(err)
			}
		case "date_range":
			start := parsed.View.State.Values["start"]["start"].SelectedDate
			end := parsed.View.State.Values["end"]["end"].SelectedDate

			if start == "" && end == "" {
				w.Header().Add("Content-Type", "application/json")
				response, _ := json.Marshal(slack.ViewSubmissionResponse{
					ResponseAction: slack.RAErrors,
					Errors: map[string]string{
						"start": "Please select either a start date or an end date.",
					},
				})
				_, err = w.Write(response)
				if err != nil {
					log.Println(err)
				}
			} else {
				tz, err := util.GetUserTimezone(parsed.User.ID)
				if err != nil {
					log.Println(err)
				}

				loc, _ := time.LoadLocation(tz)

				startDate, _ := time.ParseInLocation("2006-01-02", start, loc)
				endDate, _ := time.ParseInLocation("2006-01-02", end, loc)

				err = db.SetUserDates(startDate, endDate, parsed.User.ID)
				if err != nil {
					log.Println(err)
				}

				err = db.SetReplyMode(parsed.User.ID, db.ReplyModeDate)
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
}
