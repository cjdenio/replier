package db

import (
	"time"

	"github.com/slack-go/slack"
)

type ReplyMode string

const (
	ReplyModeManual   ReplyMode = "manual"
	ReplyModeDate     ReplyMode = "date"
	ReplyModePresence ReplyMode = "presence"
)

// User represents a DB user.
type User struct {
	Token    string    `bson:"token"`
	UserID   string    `bson:"user_id"`
	Reply    UserReply `bson:"reply"`
	Scopes   []string  `bson:"scopes"`
	TeamID   string    `bson:"team_id"`
	APIToken string    `bson:"api_token"`
}

// UserReply represents a user's chosen auto reply
type UserReply struct {
	Message   string    `bson:"message" json:"message"`
	Active    bool      `bson:"active" json:"active"`
	Whitelist []string  `bson:"whitelist" json:"whitelist"`
	Start     time.Time `bson:"start" json:"start"`
	End       time.Time `bson:"end" json:"end"`
	Mode      ReplyMode `bson:"mode" json:"mode"`
}

// ReplyShouldSend figures out whether or not the configured autoreply should be sent
func (user User) ReplyShouldSend() bool {
	if user.Reply.Message == "" {
		return false
	}

	switch user.Reply.Mode {
	case ReplyModeManual:
		return user.Reply.Active
	case ReplyModeDate:
		now := time.Now()

		if user.Reply.Start != (time.Time{}) && user.Reply.Start.After(now) {
			return false
		}

		if user.Reply.End != (time.Time{}) && user.Reply.End.Add(24*time.Hour).Before(now) {
			return false
		}
	case ReplyModePresence:
		client := slack.New(user.Token)
		if presence, err := client.GetUserPresence(user.UserID); err != nil {
			return false
		} else if presence.Presence == "active" {
			return false
		} else if presence.Presence == "away" {
			return true
		}
	}

	return true
}

// Conversation represents a single DM or channel.
type Conversation struct {
	UserID         string `bson:"user_id"`
	ConversationID string `bson:"conversation_id"`
	LastPostedOn   int64  `bson:"last_posted_on"`
}

// Installation represents an app installation
type Installation struct {
	TeamID string   `bson:"team_id"`
	Scopes []string `bson:"scopes"`
	Token  string   `bson:"token"`
	BotID  string   `bson:"bot_id"`
}
