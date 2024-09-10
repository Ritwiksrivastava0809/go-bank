package accountController

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/accounts"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (con *AccountController) CreateAccountHandler(c *gin.Context) {

	var account accounts.CreateAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		log.Error().Msgf(errorLogs.BindingJsonError, err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	if err := utils.ValidateStruct(account); err != nil {
		log.Error().Msgf("Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	if account.Owner == "" {
		log.Error().Msg("Owner field is required")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Owner field is required"})
		return
	}

	if account.Currency == "" {
		log.Error().Msg("Currency field is required")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Currency field is required"})
		return
	}

	dB := c.MustGet(constants.ConstantDB).(*db.Store)

	// Check if account already exists
	_, err := dB.GetAccountByOwner(c, account.Owner)
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
		Owner:    account.Owner,
		Currency: account.Currency,
		Balance:  100,
	}

	accounts, err := dB.CreateAccount(c, arg)

	if err != nil {
		log.Error().Msgf(errorLogs.GetCreateAccountError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Created", "ID": accounts.ID})
}

func (con *AccountController) GetAccountHandler(c *gin.Context) {
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

func (con *AccountController) UpdateAccountBalanceHandler(c *gin.Context) {
	userIDStr := c.Request.Header.Get(constants.UserID)
	if userIDStr == "" {
		log.Error().Msg("UserID header is required")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "UserID header is required"})
		return
	}

	var balance accounts.UpdateBalance
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

func (con *AccountController) AddAccountBalanaceHandler(c *gin.Context) {

	userIDStr := c.Request.Header.Get(constants.UserID)
	if userIDStr == "" {
		log.Error().Msg("UserID header is required")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "UserID header is required"})
		return
	}

	var balance accounts.UpdateBalance
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

	arg := db.AddAccountBalanceParams{
		ID:     int64(userID),
		Amount: balance.Balance,
	}

	dB := c.MustGet(constants.ConstantDB).(*db.Store)

	account, err := dB.AddAccountBalance(c, arg)
	if err != nil {
		log.Error().Msgf(errorLogs.AddAmountError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance Updated", "Account": account})
}

func (con *AccountController) ListAccountsHandler(c *gin.Context) {

	limit := c.Query(constants.PageLimit)
	offset := c.Query(constants.PageOffset)

	sortBy := c.DefaultQuery(constants.SortBy, constants.ID)
	sortDirection := c.DefaultQuery(constants.SortDirection, constants.Ascending)

	// Validate required parameters
	if limit == "" || offset == "" {
		log.Error().Msg("Limit and Offset query parameters are required")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Limit and Offset query parameters are required"})
		return
	}

	lim, err := strconv.Atoi(limit)
	if err != nil || lim <= 0 {
		log.Error().Msgf("Invalid Limit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid Limit"})
		return
	}

	off, err := strconv.Atoi(offset)
	if err != nil || off < 0 {
		log.Error().Msgf("Invalid Offset: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid Offset"})
		return
	}

	// Validate sortDirection
	sortDirection = strings.ToLower(sortDirection)
	if sortDirection != constants.Ascending && sortDirection != constants.Descending {
		log.Error().Msg("Invalid sortDirection, must be 'asc' or 'desc'")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid sortDirection, must be 'asc' or 'desc'"})
		return
	}

	// Validate sortBy
	validSortBy := map[string]bool{"id": true, "created_at": true, "balance": true, "owner": true}
	if !validSortBy[sortBy] {
		log.Error().Msg("Invalid sortBy value")
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid sortBy value"})
		return
	}

	dB := c.MustGet(constants.ConstantDB).(*db.Store)

	// Prepare query arguments for ListAccountsParams
	arg := db.ListAccountsParams{
		Limit:   int32(lim),
		Offset:  int32(off),
		Column3: sortBy,
	}

	// Execute the query
	account, err := dB.ListAccounts(c, arg)
	if err != nil {
		log.Error().Msgf(errorLogs.ListAccountsError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	// If sorting is descending, reverse the results in Go
	if sortDirection == "desc" {
		accounts.Reverse(account)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Accounts Found", "Accounts": account})
}
