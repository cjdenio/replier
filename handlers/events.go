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
		_, err := w.Write([]byte("Not verified :("))
		if err != nil {
			log.Println(err)
		}
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
		_, err = w.Write([]byte(r.Challenge))
		if err != nil {
			log.Println(err)
		}
	} else if slackEvent.Type == slackevents.CallbackEvent {
		_, err = w.Write(nil)
		if err != nil {
			log.Println(err)
		}
		innerEvent := slackEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			// If this is a message sub-event, ignore it
			if ev.SubType != "" {
				return
			}

			if ev.ChannelType == "im" {
				events.HandleMessage(slackEvent.Data.(*slackevents.EventsAPICallbackEvent), ev)
			} else {
				events.HandleMessageNonDM(slackEvent.Data.(*slackevents.EventsAPICallbackEvent), ev)
			}
		case *slackevents.AppHomeOpenedEvent:
			events.HandleAppHomeOpened(slackEvent.Data.(*slackevents.EventsAPICallbackEvent), ev)
		}
	}
}
