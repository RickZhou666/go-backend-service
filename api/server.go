package api

import (
	db "github.com/RickZhou666/go-backend-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store // remove * from server.go, as it's no longer struct pointer type, but interface instead
	router *gin.Engine
}

// NewServer creats a new HTTP server and setup routing
func NewServer(store db.Store) *Server { // remove * from server.go, as it's no longer struct pointer type, but interface instead
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	// return error msg
	return gin.H{"error": err.Error()}
}
