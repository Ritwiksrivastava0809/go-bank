package userController

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
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

	c.JSON(http.StatusOK, gin.H{"message": "user created", "User": createdUser.Username})
}
