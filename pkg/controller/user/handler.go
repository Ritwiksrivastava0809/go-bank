package userController

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"


	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/users"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)


func (con *UserController) CreateUserHandler(c *gin.Context) {


	var user users.CreateUser
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error().Msgf(errorLogs.BindingJsonError, err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	if err := utils.ValidateStruct(user); err != nil {
		log.Error().Msgf("Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	if user.Owner == "" {
		log.Error().Msg("Owner field is required")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Owner field is required"})
		return
	}

	if user.Currency == "" {
		log.Error().Msg("Currency field is required")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Currency field is required"})
		return
	}

	dB := c.MustGet(constants.ConstantDB).(*db.Store)

	// Check if account already exists
	_, err := dB.GetAccountByOwner(c, user.Owner)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Handle the case where no rows were found
			log.Info().Msg("Account not found, creating new account")
		} else {
			// Handle other errors
			log.Error().Msgf(errorLogs.GetAccountError, err)
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Error retrieving account"})
			return
		}
	} else {
		// Account exists, handle as needed (e.g., return an error or update the existing account)
		log.Info().Msg("Account already exists")
		c.JSON(http.StatusConflict, gin.H{"Message": "Account with this owner already exists"})
		return
	}

	arg := db.CreateAccountParams{
		Owner:    user.Owner,
		Currency: user.Currency,
		Balance:  0,
	}

	account, err := dB.CreateAccount(c, arg)

	if err != nil {
		log.Error().Msgf(errorLogs.GetCreateAccountError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Created", "ID": account.ID})
}

func (con *UserController) GetUserHandler(c *gin.Context) {
	dB := c.MustGet(constants.ConstantDB).(*db.Store)

	userIDStr := c.Query(constants.UserID) // Get the user ID as a string from query parameters
	if userIDStr == "" {
		log.Error().Msg("UserID query parameter is required")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "UserID query parameter is required"})
		return
	}

	// Try to convert the userID to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Error().Msgf("Invalid UserID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid UserID"})
		return
	}

	// Fetch account based on userID
	account, err := dB.GetAccount(c, int64(userID)) // Assuming GetAccount expects int64
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Handle the case where no rows were found
			log.Info().Msg("Account not found")
			c.JSON(http.StatusNotFound, gin.H{"Message": "Account not found"})
			return
		} else {
			// Handle other errors
			log.Error().Msgf(errorLogs.GetAccountError, err)
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Error retrieving account"})
			return
		}
	}

	// Successfully found the account
	c.JSON(http.StatusOK, gin.H{"message": "User Found", "Account": account})
}

func (con *UserController) UpdateUserBalanceHandler(c *gin.Context) {
	userIDStr := c.Request.Header.Get(constants.UserID)
	if userIDStr == "" {
		log.Error().Msg("UserID header is required")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "UserID header is required"})
		return
	}

	var balance users.UpdateBalance
	if err := c.ShouldBindJSON(&balance); err != nil {
		log.Error().Msgf(errorLogs.BindingJsonError, err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	if err := utils.ValidateStruct(balance); err != nil {
		log.Error().Msgf(errorLogs.ValidationError, err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Error().Msgf(errorLogs.InvalidUserError, err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid User ID format"})
		return
	}

	arg := db.UpdateAccountParams{
		ID:      int64(userID),
		Balance: balance.Balance,
	}

	dB := c.MustGet(constants.ConstantDB).(*db.Store)

	account, err := dB.UpdateAccount(c, arg)
	if err != nil {
		log.Error().Msgf(errorLogs.UpdateAccountError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance Updated", "Account": account})
}
