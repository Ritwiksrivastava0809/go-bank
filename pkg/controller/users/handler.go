package userController

import (
	"database/sql"
	"net/http"
	"net/mail"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/config"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/token"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/users"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (con *UserController) CreateUserHandler(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error().Msg(errorLogs.BindingJsonError)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorLogs.BindingJsonError})
		return
	}

	// Validate email format
	if _, err := mail.ParseAddress(user.Email); err != nil {
		log.Error().Msg(errorLogs.InvalidEmailFormatError)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}

	// Additional struct validation
	if err := utils.ValidateStruct(user); err != nil {
		log.Error().Msgf(errorLogs.ValidationError, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorLogs.ValidationError})
		return
	}

	dB := c.MustGet(constants.ConstantDB).(*db.Store)

	// Start a transaction
	tx, err := dB.DB.BeginTx(c, nil)
	if err != nil {
		log.Error().Msgf(errorLogs.TransactionError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorLogs.TransactionError})
		return
	}

	// Use the transaction with Queries
	txQueries := dB.WithTx(tx)

	// Check for existing user with the same email or username
	exists, err := txQueries.CheckExistingUser(c, db.CheckExistingUserParams{Username: user.UserName, Email: user.Email})
	if err != nil {
		tx.Rollback()
		log.Error().Msgf(errorLogs.CheckUserExistenceError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorLogs.CheckUserExistenceError})
		return
	}

	if exists {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or username already in use"})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPasswordArgon2(user.Password)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf(errorLogs.HashPasswordError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Prepare the parameters for user creation
	arg := db.CreateUserParams{
		Username:          user.UserName,
		Email:             user.Email,
		FullName:          user.FullName,
		HashedPassword:    hashedPassword,
		PasswordChangedAt: time.Now(), // Set to current time
	}

	// Create the user
	createdUser, err := txQueries.CreateUser(c, arg)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf(errorLogs.CreateUserError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Error().Msgf(errorLogs.CommitTransactionError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return the response
	resp := users.NewUserResponse(db.User{
		Username:          createdUser.Username,
		FullName:          createdUser.FullName,
		Email:             createdUser.Email,
		PasswordChangedAt: createdUser.PasswordChangedAt,
		CreatedAt:         createdUser.CreatedAt,
	})

	c.JSON(http.StatusOK, gin.H{"message": "user created", "Response": resp})
}

func (con *UserController) LoginUserHandler(c *gin.Context) {
	var login users.LoginUserRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		log.Error().Msg(errorLogs.BindingJsonError)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorLogs.BindingJsonError})
		return
	}

	dB := c.MustGet(constants.ConstantDB).(*db.Store)

	// Start a transaction
	tx, err := dB.DB.BeginTx(c, nil)
	if err != nil {
		log.Error().Msgf(errorLogs.TransactionError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorLogs.TransactionError})
		return
	}

	// Use the transaction with Queries
	txQueries := dB.WithTx(tx)

	// Get the user
	user, err := txQueries.GetUserByUsername(c, login.UserName)
	if err != nil {
		// Handle the case where the user doesn't exist
		if err == sql.ErrNoRows {
			tx.Rollback()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		// Other errors related to the query
		tx.Rollback()
		log.Error().Msgf(errorLogs.GetUserError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorLogs.GetUserError})
		return
	}

	// Check if the password is correct
	if err := utils.VerifyPassword(user.HashedPassword, login.Password); err != nil {
		tx.Rollback()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Error().Msgf(errorLogs.CommitTransactionError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Create token and send it back
	token, ok := c.MustGet(constants.TokenMaker).(token.Maker)
	if !ok {
		log.Error().Msg("failed to retrieve token maker from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	accessToken, err := token.CreateToken(
		user.Username,
		time.Duration(config.GetAccessTokenDuration()),
	)

	if err != nil {
		log.Error().Msgf(errorLogs.TokenError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := users.LoginUserResponse{
		AccessToken: accessToken,
		User:        users.NewUserResponse(user),
	}

	// Add response if needed (e.g., JWT token generation or success message)
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "response": resp})
}
