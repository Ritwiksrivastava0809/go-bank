package server

import (
	"fmt"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/config"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/rs/zerolog/log"
)

func Init(dbConnection *db.Store) {
	config := config.GetConfig()
	router := NewServer(dbConnection)
	log.Debug().Msg(config.GetString("server.port"))
	serverUrl := fmt.Sprintf("%s:%s", config.GetString("server.host"), config.GetString("server.port"))
	router.router.Run(serverUrl)
}
