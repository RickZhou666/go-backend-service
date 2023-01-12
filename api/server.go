package api

import (
	"fmt"

	db "github.com/RickZhou666/go-backend-service/db/sqlc"
	"github.com/RickZhou666/go-backend-service/token"
	"github.com/RickZhou666/go-backend-service/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store // remove * from server.go, as it's no longer struct pointer type, but interface instead
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creats a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) { // remove * from server.go, as it's no longer struct pointer type, but interface instead
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

// split rounter
func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	// return error msg
	return gin.H{"error": err.Error()}
}
