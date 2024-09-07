package db

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// DB connection details
const (
	DBDriver = "postgres"
	DBSource = "postgresql://root:password@localhost:5432/bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(DBDriver, DBSource)
	if err != nil {
		message := fmt.Sprintf("got error getting db %v", err)
		log.Error().Msg("Got error while getting db in test main error :: " + message)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
