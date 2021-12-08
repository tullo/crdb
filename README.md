# crdb

- Go [ORM](company/readme.md) libs

## Integration Test

Launches a single node cockroachdb in a docker container using the [github.com/ory/dockertest/v3/docker](https://github.com/ory/dockertest) package and executes the test function.

```sh
# 1. Pulls docker image cockroachdb/cockroach:v21.2.2
# 2. Executes tests located in roach_test.go
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
```

## Migrations Test

Launches a single node cockroachdb in a docker container using the [github.com/golang-migrate/migrate/v4/dktesting](https://github.com/dhui/dktest) package and is applying all up migrations when executing the test function using [github.com/golang-migrate/migrate/v4](github.com/golang-migrate/migrate).

```sh
# 1. Pulls docker image cockroachdb/cockroach:v21.2.2
# 2. Processes migration files (18) located under db/migrations
#    by applying all up migrations.
go test -v -timeout 30s -run ^TestMigrate$ github.com/tullo/crdb

=== RUN   TestMigrate
=== RUN   TestMigrate/cockroachdb/cockroach:v21.2.2
=== PAUSE TestMigrate/cockroachdb/cockroach:v21.2.2
=== CONT  TestMigrate/cockroachdb/cockroach:v21.2.2
    dktest.go:36: Pulling image: cockroachdb/cockroach:v21.2.2        
    dktest.go:80: Created container: ...
    dktest.go:85: Started container: ...
    dktest.go:95: Inspected container: ...
    dktest.go:130: Stopped container: ...
    dktest.go:136: Removed container: ...
--- PASS: TestMigrate (0.00s)
    --- PASS: TestMigrate/cockroachdb/cockroach:v21.2.2 (7.52s)
PASS
ok  	github.com/tullo/crdb	7.527s
```
