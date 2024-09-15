package transactionController

import (
	"net/http"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/token"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/transactions"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (con *TransactionController) InsertTransactionHandler(c *gin.Context) {

	var txn transactions.Transaction
	if err := c.ShouldBindJSON(&txn); err != nil {
		log.Error().Msgf(errorLogs.BindingJsonError, err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	fromAccount, valid := con.ValidAccount(c, txn.FromAccountID, txn.Currency)
	if !valid {
		return
	}

	authPayload := c.MustGet(constants.AuthorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != authPayload.Username {
		log.Error().Msg("Error :: Account does not belong to the user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account does not belong to the user"})
		return
	}

	_, valid = con.ValidAccount(c, txn.ToAccountID, txn.Currency)
	if !valid {
		return
	}

	if err := utils.ValidateStruct(txn); err != nil {
		log.Error().Msgf(errorLogs.ValidationError, err)
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: txn.FromAccountID,
		ToAccountID:   txn.ToAccountID,
		Amount:        txn.Amount,
	}

	dB := c.MustGet(constants.ConstantDB).(*db.Store)

	transferTxResult, err := dB.TransferTx(c, arg)
	if err != nil {
		log.Error().Msgf(errorLogs.TransferTxError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Insert Transaction", "data": transferTxResult})
}
