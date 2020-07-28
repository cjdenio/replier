package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/cjdenio/replier/db"
	"github.com/cjdenio/replier/util"
	"github.com/slack-go/slack"
)

// HandleInteractivity handles interactions in Slack
func HandleInteractivity(w http.ResponseWriter, r *http.Request) {
	buf, _ := ioutil.ReadAll(r.Body)
	r.Form, _ = url.ParseQuery(string(buf))
	//fmt.Println(r.Form.Get("payload"))

	var parsed slack.InteractionCallback
	json.Unmarshal([]byte(r.Form.Get("payload")), &parsed)

	if parsed.Type == slack.InteractionTypeBlockActions {
		w.Write(nil)

		switch parsed.ActionCallback.BlockActions[0].ActionID {
		case "edit_message":
			user, _ := db.GetUser(parsed.User.ID)

			client := slack.New(os.Getenv("SLACK_TOKEN"))
			_, err := client.OpenView(parsed.TriggerID, slack.ModalViewRequest{
				Type:       "modal",
				Title:      slack.NewTextBlockObject("plain_text", "Edit Message", false, false),
				CallbackID: "edit_message",
				Blocks: slack.Blocks{
					BlockSet: []slack.Block{
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
					},
				},
				Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
				Submit: slack.NewTextBlockObject("plain_text", "Save", false, false),
			})

			if err != nil {
				log.Fatal(err)
			}
		case "reply_toggle":
			db.ToggleReplyActive(parsed.User.ID)
			util.UpdateAppHome(parsed.User.ID)
		}
	} else if parsed.Type == slack.InteractionTypeViewSubmission {
		w.Write(nil)

		switch parsed.View.CallbackID {
		case "edit_message":
			desiredMessage := parsed.View.State.Values["message"]["message"].Value

			db.SetUserMessage(parsed.User.ID, desiredMessage)
			util.UpdateAppHome(parsed.User.ID)
		}
	}
}
