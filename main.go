package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/config"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/logger"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/server"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/rs/zerolog/log"
)

func main() {

	logger.InitLogger()

	environment := flag.String("e", "development", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()

	config.Init(*environment)

	utils.InitValidator()

	database, err := config.NewDB()
	if err != nil {
		log.Fatal().Msg("Cannot connect to database")
	}
	defer database.Close()

	store := db.NewStore(database)
	server.Init(store)
	log.Info().Msg("Server started")
}
