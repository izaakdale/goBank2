package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/izaakdale/goBank2/db/sqlc"
)

// serves all http requests
type Server struct {
	store  db.Store
	router *gin.Engine
}

// creates a new http server and sets up routing
func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.CreateTransfer)

	server.router = router
	return server
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
