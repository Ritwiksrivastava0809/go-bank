package server

import (
	"net/http"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	userController "github.com/Ritwiksrivastava0809/go-bank/pkg/controller/user"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/middleware"

	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server serves htttp requests
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set(constants.ConstantDB, store)
		c.Next()
	})

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(middleware.LoggerMiddleware())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}
	corsConfig.AllowHeaders = []string{constants.Origin, constants.ContentType, constants.ContentLength, constants.Authorization}

	// Apply the CORS middleware to your router
	router.Use(cors.New(corsConfig))

	// Initialize the routes
	v0 := router.Group("/v0")
	{
		userGroup := v0.Group("/users")
		{
			userController := new(userController.UserController)
			userGroup.POST("/create", userController.CreateUser)
		}
	}

	server.router = router
	return server

}
