package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/config"
	db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/logger"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/server"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
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

	// Initialize New Relic application
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("go-bank"),
		newrelic.ConfigLicense("8e816dd73b01fe07b7376bed16bbb0acFFFFNRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogDecoratingEnabled(true),
	)

	if err != nil {
		logrus.Fatalf("New Relic initialization failed: %s", err)
	}

	// apikey := "8e816dd73b01fe07b7376bed16bbb0acFFFFNRAL"

	// Initialize database connection
	database, err := config.NewDB()
	if err != nil {
		logrus.Fatalf("Cannot connect to database error :: %s", err)
	}
	defer database.Close()

	store := db.NewStore(database)

	// Initialize the server with New Relic instrumentation
	server.InitWithNewRelic(store, app)

	logrus.Info("Server started")

	// Ensure New Relic application shutdown
	defer app.Shutdown(10 * time.Second)
}
