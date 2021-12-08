# Company API using GORM

- `go get gorm.io/gorm@latest`
- `go get gorm.io/driver/postgres@latest`
- DB [drivers](https://github.com/orgs/go-gorm/repositories?q=driver&type=&language=&sort=)

## TestCompanyAPI

```sh
# launch crdb instance.
cd docker
docker-compose --verbose up -d

# connect and create test database.
cockroach sql --insecure --url postgresql://0.0.0.0:26257
root@0.0.0.0:26257/defaultdb> CREATE DATABASE IF NOT EXISTS company_gorm;

# build and start the app.
cd ../company/gorm
go build && ./gorm

# run end-to-end tests
go test -v github.com/tullo/crdb/company/gorm
=== RUN   TestCompanyAPI
    api_test.go:55: createCustomer: 717324675344072705
    api_test.go:77: getProducts: 1
    api_test.go:124: createOrder: 717324675419406337
    api_test.go:34: Order &{ID:717324675419406337 Subtotal:630.08 Customer:{ID:717324675344072705 Name:0xc000294250} CustomerID:0 Products:[{ID:717320840834842625 Name:0xc000294260 Price:315.04}]}
--- PASS: TestCompanyAPI (0.05s)
PASS
ok  	github.com/tullo/crdb/company/gorm	0.056s

# cleaning up
root@0.0.0.0:26257/defaultdb> DROP DATABASE company_gorm CASCADE;

cd ../../docker
docker-compose down -v
```
