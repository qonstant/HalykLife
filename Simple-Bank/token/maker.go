package token

import (
	db "Simple-Bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"time"
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration, userRole string) (string, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string, store db.Store, ctx *gin.Context) (*Payload, error)
}
