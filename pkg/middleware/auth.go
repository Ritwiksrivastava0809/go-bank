package middleware

import (
	"net/http"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/config"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AuthInternalTokenMiddleware(c *gin.Context) {
	//get the internal token from the header
	internalToken := c.GetHeader(constants.InternalToken)
	// check if the internal token is empty
	if internalToken == "" {
		log.Error().Msg("Error no Internal Token present")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id token not available"})
		c.Abort()
		return
	}

	if internalToken != config.GetInternalToken() {
		log.Error().Msg(" Error :: Incorrect Internal Token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}
	c.Next()
}
