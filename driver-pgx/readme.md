# PostgreSQL driver and toolkit for Go

- https://github.com/jackc/pgx
- https://pkg.go.dev/github.com/jackc/pgx/v5

## Driver

`go get github.com/jackc/pgx/v5@latest`

## Bank Example

[Covered steps](https://github.com/cockroachlabs/hello-world-go-pgx):
- Create a table in the bank database.
- Insert some rows into the table you created.
- Read values from the table.
- Execute a batch of statements as an atomic [transaction](https://www.cockroachlabs.com/docs/v20.2/transactions).

Note that CockroachDB may require the client to retry a transaction in the case of read/write contention. The CockroachDB [Go client](https://github.com/cockroachdb/cockroach-go) includes a generic **retry function** `(ExecuteTx())` that runs inside a transaction and retries it as needed. The code sample shows how you can use this function to wrap SQL statements.

```sql
CREATE DATABASE bank;
CREATE USER johndoe WITH PASSWORD gopher;
GRANT ALL ON DATABASE bank TO johndoe;
-- DROP DATABASE bank CASCADE;
```

`go build .`

`./driver-pgx`

```console
Initial balances:
1 1000
2 250
Success
```

## DB

```sh
cd docker
docker-compose up -d

cockroach sql --insecure --url postgresql://0.0.0.0:26257
```
