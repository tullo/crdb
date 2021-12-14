# Bun

https://bun.uptrace.dev/

> Bun comes with its own PostgreSQL driver called pgdriver.

# Company API using bun/pgdriver

- `go get github.com/uptrace/bun/driver/pgdriver@v1.0.19`

## TestCompanyAPI

```sh
# launch an insecure crdb instance for testing.
cd ../../ && make

# connect and create test database.
cockroach sql --insecure -e "CREATE DATABASE IF NOT EXISTS company_bun;"

# build and start the app.
cd ../company/gopg
go build && ./gopg

# run end-to-end tests
go test -v github.com/tullo/crdb/company/gobun -count=1
=== RUN   TestCompanyAPI
    api_test.go:17: Customer {718742717374922753 John Doe}
    api_test.go:29: Product &{ID:718742717436329985 Name:GopherCon Europe 2021 Price:315.04}
    api_test.go:36: Order &{ID:718742717499080705 Subtotal:630.08 Customer:{ID:718742717374922753 Name:John Doe} CustomerID:0 Products:[{ID:718742717436329985 Name:GopherCon Europe 2021 Price:315.04}]}
--- PASS: TestCompanyAPI (0.11s)
PASS
ok  	github.com/tullo/crdb/company/gobun	0.115s

# cleaning up
cockroach sql --insecure -e "DROP DATABASE company_bun CASCADE;"
cd ../../ && make clean
```

## Bun realworld application

https://github.com/tullo/bun-realworld-app
