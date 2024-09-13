package token

import (
	"fmt"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	"github.com/golang-jwt/jwt"
)

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < constants.MinSecretKeyLen {
		return nil, fmt.Errorf(errorLogs.InvalidKeySize, constants.MinSecretKeyLen)
	}

	return &JWTMAKER{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMAKER) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMAKER) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf(constants.InvalidTokenError)
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && verr.Errors == jwt.ValidationErrorExpired {
			return nil, fmt.Errorf(constants.ExipredToken)
		}
		return nil, fmt.Errorf(constants.InvalidTokenError)
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf(constants.InvalidTokenError)
	}

	return payload, nil
}
