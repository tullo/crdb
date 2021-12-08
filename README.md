# crdb

## Integration Test

Launches a single node cockroachdb in a docker container using the `"github.com/ory/dockertest/v3/docker"` package and executes the test function.

```sh
# 1. pulls image cockroachdb/cockroach:v21.2.2
# 2. executes test located in roach_test.go
go test -v -timeout 30s -run ^TestRoachIntegration$ github.com/tullo/crdb

=== RUN   TestRoachIntegration
    roach_test.go:58: ping postgresql://admin@0.0.0.0:26257?sslmode=disable
    roach_test.go:58: ping postgresql://admin@0.0.0.0:26257?sslmode=disable
    roach_test.go:73: Stats {MaxOpenConnections:0 OpenConnections:1 InUse:0 Idle:1 WaitCount:0 WaitDuration:0s MaxIdleClosed:0 MaxIdleTimeClosed:0 MaxLifetimeClosed:0}
    roach_test.go:140: Balances: 
    roach_test.go:146: 1 1000
    roach_test.go:146: 2 250
    roach_test.go:140: Balances: after tx
    roach_test.go:146: 1 900
    roach_test.go:146: 2 350
--- PASS: TestRoachIntegration (1.55s)
PASS
ok  	github.com/tullo/crdb	1.560s

````