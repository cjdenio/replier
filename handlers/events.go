package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cjdenio/replier/handlers/events"
	"github.com/cjdenio/replier/util"
	"github.com/slack-go/slack/slackevents"
)

// HandleEvents handles Events API requests
func HandleEvents(w http.ResponseWriter, r *http.Request) {
	buf, _ := ioutil.ReadAll(r.Body)

	if !util.VerifySlackRequest(r, buf) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not verified :("))
		return
	}

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
			if ev.ChannelType == "im" {
				//fmt.Println(string(buf))
				events.HandleMessage(slackEvent.Data.(*slackevents.EventsAPICallbackEvent), ev)
			}
		case *slackevents.AppHomeOpenedEvent:
			events.HandleAppHomeOpened(slackEvent.Data.(*slackevents.EventsAPICallbackEvent), ev)
		}
	}
}
