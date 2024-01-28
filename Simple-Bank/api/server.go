package api

import (
	db "Simple-Bank/db/sqlc"
	"Simple-Bank/token"
	"Simple-Bank/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.setupRouter()

	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker, server.store))

	authRoutes.GET("/accounts", server.listUsers)
	authRoutes.GET("/user/:iin", server.getQrForUser)

	//authRoutes.GET("/accounts/:id", server.getAccount)
	//authRoutes.GET("/accounts", server.listAccounts)
	//
	//authRoutes.POST("/transfers", server.createTransfer)
	//authRoutes.GET("/transfers", server.listTransfers)

	server.router = router
}
