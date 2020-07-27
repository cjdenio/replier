package db

// User represents a DB user.
type User struct {
	Token  string    `json:"token"`
	UserID string    `json:"user_id"`
	Reply  UserReply `json:"reply"`
}

// UserReply represents a user's chosen auto reply
type UserReply struct {
	Message string `json:"message"`
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Active  bool   `json:"active"`
}

// ReplyShouldSend figures out whether or not the configured autoreply should be sent
func (user User) ReplyShouldSend() bool {
	if user.Reply.Message == "" || !user.Reply.Active {
		return false
	}

	return true
}
