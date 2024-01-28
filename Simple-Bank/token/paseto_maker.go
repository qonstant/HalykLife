package token

import (
	db "Simple-Bank/db/sqlc"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d chars, urs is %d, which is: %s", chacha20poly1305.KeySize, len(symmetricKey), symmetricKey)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration, userRole string) (string, error) {
	payload, err := NewPayload(username, duration, userRole)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string, store db.Store, ctx *gin.Context) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	role, err := store.GetRoleByUsername(ctx, payload.Username)

	if err != nil {
		return nil, err
	}
	if role != payload.UserRole {
		return nil, err
	}

	return payload, nil
}
