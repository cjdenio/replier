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
}

// HasActiveReply figures out whether or not the user has a currently active reply
func (user User) HasActiveReply() bool {
	if user.Reply.Message == "" {
		return false
	}

	return true
}
