# Testserver

The testserver package helps running CRDB binary with tests.

```go
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
```

## Test

```sh
go test -race -v

=== RUN   TestRunServer
2022/07/06 10:00:05 GET https://binaries.cockroachdb.com/cockroach-v22.1.2.linux-amd64.tgz
2022/07/06 10:00:05 Using automatically-downloaded binary: /tmp/cockroach-22-1-2
2022/07/06 10:00:05 executing: /tmp/cockroach-22-1-2 start-single-node --logtostderr --insecure --host=localhost --port=0 --http-port=0 --store=type=mem,size=0.20 --listening-url-file=/tmp/cockroach-testserver3503377092/listen-url0
2022/07/06 10:00:05 process 93304 started: /tmp/cockroach-22-1-2 start-single-node --logtostderr --insecure --host=localhost --port=0 --http-port=0 --store=type=mem,size=0.20 --listening-url-file=/tmp/cockroach-testserver3503377092/listen-url0
    roach_test.go:178: Balances: 
    roach_test.go:184: 1 1000
    roach_test.go:184: 2 250
    roach_test.go:178: Balances: after tx
    roach_test.go:184: 1 900
    roach_test.go:184: 2 350
--- PASS: TestWithTestServer (1.18s)
PASS
ok  	github.com/tullo/crdb	1.213s
```
