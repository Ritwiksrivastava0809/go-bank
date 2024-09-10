package transactionController

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (con *TransactionController) ValidAccount(c *gin.Context, accountID int64, currency string) bool {
	dB := c.MustGet(constants.ConstantDB).(*db.Store)
	flag := true

	account, err := dB.GetAccount(c, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Handle the case where no rows were found
			log.Info().Msg("Account not found, creating new account")
			return false
		}
		// Handle other errors
		log.Error().Msgf(errorLogs.GetAccountError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Error retrieving account"})
		return false

	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s vs %s", accountID, account.Currency, currency)
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return false
	}

	return flag
}
