package db

// User represents a DB user.
type User struct {
	Token  string    `json:"token"`
	UserID string    `json:"user_id"`
	Reply  UserReply `json:"reply"`
}

// UserReply represents a user's chosen auto reply
type UserReply struct {
	Message   string   `json:"message"`
	Active    bool     `json:"active"`
	Whitelist []string `json:"whitelist"`
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
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	LastPostedOn   int    `json:"last_posted_on"`
}
