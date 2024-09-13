package token

import (
	"time"

	"github.com/google/uuid"
)

// Maker is an interface that wraps the Sign and Verify methods
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// JWTMAKER is a struct that contains the secret key
type JWTMAKER struct {
	secretKey string
}
