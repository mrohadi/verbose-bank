package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/mrohadi/simplebank/db/sqlc"
	"github.com/mrohadi/simplebank/token"
	"github.com/mrohadi/simplebank/utils"
)

// Server serve HTTP request for banking system
type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer create new HTTP server and routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmectricKey)
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

// setupRouter()
func (s *Server) setupRouter() {
	router := gin.Default()

	// users routing
	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)

	// accounts routing
	router.POST("/accounts", s.createAccount)
	router.GET("/accounts/:id", s.getAccount)
	router.GET("/accounts", s.listAccount)

	// transfer routing
	router.POST("/transfers", s.createTransfer)

	s.router = router
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(addr string) error {
	return s.router.Run()
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
