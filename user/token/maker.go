package token

import "time"

// Maker is an interface for managing tokens.
type Maker interface {
	// CreateToken creates a new token for a specific user id, username and duration
	CreateToken(userID int64, username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
