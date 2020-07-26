package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slack-go/slack/slackevents"
)

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	buf, _ := ioutil.ReadAll(r.Body)

	slackEvent, err := slackevents.ParseEvent(buf, slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Fatal(err)
	}

	if slackEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(buf, &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	} else if slackEvent.Type == slackevents.CallbackEvent {
		w.Write([]byte("cool"))
		innerEvent := slackEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			fmt.Println([]byte(ev.Text))
		}
	}
}
