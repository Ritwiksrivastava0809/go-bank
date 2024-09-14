package token

import (
	"fmt"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	Payload := &Payload{
		ID:        token,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return Payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return fmt.Errorf(constants.ExipredToken)
	}

	return nil
}

// PasetoMaker is a Paseto token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}
