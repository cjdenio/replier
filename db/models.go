package db

// User represents a DB user.
type User struct {
	Token  string    `bson:"token"`
	UserID string    `bson:"user_id"`
	Reply  UserReply `bson:"reply"`
}

// UserReply represents a user's chosen auto reply
type UserReply struct {
	Message   string   `bson:"message"`
	Active    bool     `bson:"active"`
	Whitelist []string `bson:"whitelist"`
}

// ReplyShouldSend figures out whether or not the configured autoreply should be sent
func (user User) ReplyShouldSend() bool {
	if user.Reply.Message == "" || !user.Reply.Active {
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
