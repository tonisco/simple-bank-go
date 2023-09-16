package token

import "time"

// Make is a unique interface for managing tokens
type Maker interface {
	// create token for a specific user name and duration
	createToken(username string, duration time.Duration) (string, error)

	// checks if token is valid
	verifyToken(token string) (*Payload, error)
}
