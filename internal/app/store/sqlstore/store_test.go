package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")

	if len(databaseURL) == 0 {
		databaseURL = "host=localhost dbname=restapi_test user=postgres password=asdadm443 sslmode=disable"
	}

	os.Exit(m.Run())
}
