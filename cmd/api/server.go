package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/mrohadi/simplebank/db/sqlc"
)

// Server serve HTTP request for banking system
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer create new HTTP server and routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// add http routing
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(addr string) error {
	return s.router.Run()
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
