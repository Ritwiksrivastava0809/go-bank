package userController

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/users"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (con *UserController) CreateUser(c *gin.Context) {

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
