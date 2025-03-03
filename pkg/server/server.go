package server

import (
	"fmt"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/config"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

func InitWithNewRelic(dbConnection *db.Store, app *newrelic.Application) {
	config := config.GetConfig()
	router, err := NewServer(dbConnection, app)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("error while initializing the server :: %s", err))
	}
	log.Debug().Msg(config.GetString("server.port"))
	serverUrl := fmt.Sprintf("%s:%s", config.GetString("server.host"), config.GetString("server.port"))
	router.router.Run(serverUrl)
}
