package token

import (
	"fmt"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// NewPasstoMaker creates a new PasetoMaker
func NewPasetoMaker(symetricKey string) (Maker, error) {

	if len(symetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf(errorLogs.InvalidKeySize, chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		symetricKey: []byte(symetricKey),
		paseto:      paseto.NewV2(),
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symetricKey, payload, nil)
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symetricKey, payload, nil)
	if err != nil {
		return nil, fmt.Errorf(constants.InvalidToken)
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
