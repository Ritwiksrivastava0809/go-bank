package db

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	DBDriver = "postgres"
	DBSource = "postgresql://root:password@localhost:5432/bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		message := fmt.Sprintf("got error getting db %v", err)
		log.Error().Msg("Got error while getting db in test main error :: " + message)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
