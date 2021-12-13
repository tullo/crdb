# PostgreSQL client and ORM for Go

https://pg.uptrace.dev/

> go-pg is in a maintenance mode and will NOT receive new features. Please use [Bun](https://bun.uptrace.dev/) instead.

# Company API using go-pg

- `go get github.com/go-pg/pg/v10@latest`

## TestCompanyAPI

```sh
# launch crdb instance.
cd docker
docker-compose --verbose up -d

# connect and create test database.
cockroach sql --insecure --url postgresql://0.0.0.0:26257
root@0.0.0.0:26257/defaultdb> CREATE DATABASE IF NOT EXISTS company_gopg;

# build and start the app.
cd ../company/gopg
go build && ./gopg

# run end-to-end tests
go test -v github.com/tullo/crdb/company/gopg
=== RUN   TestCompanyAPI
    api_test.go:55: createCustomer 717352417490436097
    api_test.go:77: getProducts 1
    api_test.go:124: createOrder 717352417528741889
    api_test.go:34: Order &{ID:717352417528741889 Subtotal:630.08 Customer:{ID:717352417490436097 Name:John Doe} CustomerID:0 Products:[{ID:717351113135816705 Name:GopherCon Europe 2021 Price:315.04}]}
--- PASS: TestCompanyAPI (0.04s)
PASS
ok  	github.com/tullo/crdb/company/gopg	0.047s

# cleaning up
root@0.0.0.0:26257/defaultdb> DROP DATABASE company_gopg CASCADE;

cd ../../docker
docker-compose down -v
```
