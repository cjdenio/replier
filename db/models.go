package db

import (
	"time"
)

// User represents a DB user.
type User struct {
	Token  string    `bson:"token"`
	UserID string    `bson:"user_id"`
	Reply  UserReply `bson:"reply,omitempty"`
	Scopes []string  `bson:"scopes"`
}

// UserReply represents a user's chosen auto reply
type UserReply struct {
	Message   string    `bson:"message"`
	Active    bool      `bson:"active"`
	Whitelist []string  `bson:"whitelist"`
	Start     time.Time `bson:"start"`
	End       time.Time `bson:"end"`
}

// ReplyShouldSend figures out whether or not the configured autoreply should be sent
func (user User) ReplyShouldSend() bool {
	if user.Reply.Message == "" || !user.Reply.Active {
		return false
	}

	now := time.Now()

	if user.Reply.Start != (time.Time{}) && user.Reply.Start.After(now) {
		return false
	}

	if user.Reply.End != (time.Time{}) && user.Reply.End.Add(24*time.Hour).Before(now) {
		return false
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
}
