package testsupport

import (
	"database/sql"
	"testing"

	"github.com/cockroachdb/cockroach-go/v2/testserver"
)

func TestRunServer(t *testing.T) {
	ts, err := testserver.NewTestServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Stop()

	db, err := sql.Open("postgres", ts.PGURL().String())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec("SELECT 1;"); err != nil {
		t.Log(err)
		t.Fail()
	}
}
