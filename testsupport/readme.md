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
2021/12/11 10:57:36 GET https://binaries.cockroachdb.com/cockroach-v21.2.2.linux-amd64.tgz
2021/12/11 10:57:36 Using automatically-downloaded binary: /tmp/cockroach-21-2-2
2021/12/11 10:57:36 executing: /tmp/cockroach-21-2-2 start-single-node --logtostderr --insecure --host=localhost --port=0 --http-port=0 --store=type=mem,size=0.20 --listening-url-file=/tmp/cockroach-testserver832636131/listen-url
2021/12/11 10:57:36 process 412373 started: /tmp/cockroach-21-2-2 start-single-node --logtostderr --insecure --host=localhost --port=0 --http-port=0 --store=type=mem,size=0.20 --listening-url-file=/tmp/cockroach-testserver832636131/listen-url
--- PASS: TestRunServer (0.75s)
PASS
ok  	github.com/tullo/crdb/testsupport	0.756s
```
