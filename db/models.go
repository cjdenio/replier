package db

// User represents a DB user.
type User struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}
