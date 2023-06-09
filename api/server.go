package api

import (
	"fmt"
	db "go-backend/db/sqlc"
	"go-backend/token"
	"go-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// The Server type contains a database store and a router for handling HTTP requests in a Go
// application.
// @property store - `store` is a pointer to an instance of `db.Store`, which is likely a struct that
// handles interactions with a database. It could contain methods for querying, inserting, updating,
// and deleting data from the database. The `Server` struct likely uses this `store` property to
// interact with the
// @property router - The `router` property is a pointer to an instance of the `gin.Engine` struct,
// which is a popular HTTP web framework for building RESTful APIs in Go. It provides a set of methods
// for defining routes, handling requests, and rendering responses. The `router` is responsible for
// mapping incoming
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// The `Start` function is a method of the `Server` struct that starts the server by running the router
// on a specified address. It takes in an `address` string parameter and returns an error if there is
// any issue starting the server. The `server.router.Run(address)` method is called to start the server
// on the specified address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// The function creates a new server instance with a given database store and sets up a router with
// routes.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	// register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// routes
	apiRouter := router.Group("/api/v1")
	server.addUserRoutes(apiRouter)
	server.addTokenRoutes(apiRouter)

	// auth routes
	apiRouter.Use(authMiddleware(server.tokenMaker))
	server.addAccountRoutes(apiRouter)
	server.addTransferRoutes(apiRouter)

	server.router = router
	return server, nil
}
