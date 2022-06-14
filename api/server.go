package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/izaakdale/goBank2/db/sqlc"
	"github.com/izaakdale/goBank2/token"
	"github.com/izaakdale/goBank2/util"
)

// serves all http requests
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

// creates a new http server and sets up setupRouting
func NewServer(config util.Config, store db.Store) (*Server, error) {

	pasetoMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create Maker, server initialization failed")
	}

	server := &Server{
		store:      store,
		tokenMaker: pasetoMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouting()
	return server, nil
}

func (server *Server) setupRouting() {

	router := gin.Default()
	// no middleware required
	router.POST("/users", server.createUser)
	router.GET("/login", server.login)
	// auth required
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	authRoutes.GET("/users/:username", server.getUser)

	server.router = router
}

// runs the server at the specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// add type map[string]interface{} to utils-go
func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
