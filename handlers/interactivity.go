package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/slack-go/slack"
)

func HandleInteractivity(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{})

	buf, _ := ioutil.ReadAll(r.Body)
	r.Form, _ = url.ParseQuery(string(buf))
	//fmt.Println(r.Form.Get("payload"))

	var parsed slack.InteractionCallback
	json.Unmarshal([]byte(r.Form.Get("payload")), &parsed)

	switch parsed.ActionCallback.BlockActions[0].ActionID {
	case "edit_message":
		client := slack.New(os.Getenv("SLACK_TOKEN"))
		_, err := client.OpenView(parsed.TriggerID, slack.ModalViewRequest{
			Type:  "modal",
			Title: slack.NewTextBlockObject("plain_text", "Edit", false, false),
			Blocks: slack.Blocks{
				BlockSet: []slack.Block{
					slack.NewSectionBlock(
						slack.NewTextBlockObject("mrkdwn", "howdy", false, false),
						nil,
						nil,
					),
				},
			},
			Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
			Submit: slack.NewTextBlockObject("plain_text", "Save", false, false),
		})

		if err != nil {
			log.Fatal(err)
		}
	}
}
