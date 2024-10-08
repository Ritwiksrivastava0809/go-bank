package server

import (
	"fmt"
	"net/http"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/config"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	accountController "github.com/Ritwiksrivastava0809/go-bank/pkg/controller/accounts"
	transactionController "github.com/Ritwiksrivastava0809/go-bank/pkg/controller/transactions"
	userController "github.com/Ritwiksrivastava0809/go-bank/pkg/controller/users"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/middleware"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/token"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/go-playground/validator/v10"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Server serves htttp requests
type Server struct {
	store     *db.Store
	tokeMaker token.Maker
	router    *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(store *db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.GetSymmetricKey())
	if err != nil {
		return nil, fmt.Errorf(errorLogs.TokenError, err)
	}

	server := &Server{
		store:     store,
		tokeMaker: tokenMaker,
	}
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set(constants.ConstantDB, store)
		c.Set(constants.TokenMaker, tokenMaker)
		c.Next()
	})

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", utils.ValidCurrency)
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(middleware.LoggerMiddleware())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}

	corsConfig.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}
	corsConfig.AllowHeaders = []string{constants.Origin, constants.ContentType, constants.ContentLength, constants.Authorization}

	// Apply the CORS middleware to your router
	router.Use(cors.New(corsConfig))

	authMiddleWare := middleware.AuthTokenMiddleware(tokenMaker)

	// Initialize the routes
	v0 := router.Group("/v0")

	{

		userGroup := v0.Group("/users")
		userGroup.Use(middleware.AuthInternalTokenMiddleware)
		{
			userController := new(userController.UserController)
			userGroup.POST("/create", userController.CreateUserHandler)
			userGroup.POST("/login", userController.LoginUserHandler)
		}

		accountGroup := v0.Group("/accounts")
		{
			accountController := new(accountController.AccountController)
			accountGroup.POST("/create", authMiddleWare, middleware.AuthInternalTokenMiddleware, accountController.CreateAccountHandler)
			accountGroup.GET("/get", authMiddleWare, middleware.AuthInternalTokenMiddleware, accountController.GetAccountHandler)
			accountGroup.PATCH("/update", middleware.AuthInternalTokenMiddleware, accountController.UpdateAccountBalanceHandler)
			accountGroup.PATCH("/add", middleware.AuthInternalTokenMiddleware, accountController.AddAccountBalanaceHandler)
			accountGroup.GET("/list", authMiddleWare, middleware.AuthInternalTokenMiddleware, accountController.ListAccountsHandler)
		}

		transactionGroup := v0.Group("/transactions")
		transactionGroup.Use(authMiddleWare)
		transactionGroup.Use(middleware.AuthInternalTokenMiddleware)
		{
			transactionController := new(transactionController.TransactionController)
			transactionGroup.POST("/Insert", transactionController.InsertTransactionHandler)
		}
	}

	server.router = router
	return server, nil

}
