# upper/db - data access layer for Go

Agnostic API - SQL friendly - ORM-like layer

- https://upper.io/v4/
- https://upper.io/v4/adapter/cockroachdb/
- https://tour.upper.io/

## Driver

`go get github.com/upper/db/v4@latest`

## Adapter

`go get github.com/upper/db/v4/adapter/cockroachdb@latest`

## Bank Example

[Covered steps](https://github.com/cockroachlabs/hello-world-go-upperdb):
- Use  `upper/db` to map Go-specific objects to SQL operations.
- Creates the accounts table, if it does not already exist.
- Deletes any existing rows in the accounts table.
- Inserts two rows into the accounts table.
- Prints the rows in the accounts table to the terminal.
- Deletes the first row in the accounts table.
- Updates the rows in the accounts table within an explicit [transaction](https://www.cockroachlabs.com/docs/v20.2/transactions).
- Simulates transaction retries triggered by `crdb_internal.force_retry`.
- Prints the rows in the accounts table to the terminal once more.

Note that CockroachDB may require the client to retry a transaction in the case of read/write contention.

```sql
CREATE DATABASE bank;
CREATE USER johndoe WITH PASSWORD gopher;
GRANT ALL ON DATABASE bank TO johndoe;
-- DROP DATABASE bank CASCADE;
```

## Upper DB

```go
import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/cockroachdb"
)
```

`go build .`

```sh
./upper

2021/12/08 21:39:58 Balances:
	accounts[717375400172322817]: 1000
	accounts[717375400193490945]: 250

2021/12/08 21:39:58 Balances:
	accounts[717375400172322817]: 500
	accounts[717375400193490945]: 999

# Simulated TX retries for 50 milliseconds
2021/12/08 21:39:58 upper/db: log_level=WARNING file=/home/anda/go/pkg/mod/github.com/upper/db/v4@v4.2.1/internal/sqladapter/session.go:646
	Session ID:     00006
	Transaction ID: 00005
	Query:          SELECT crdb_internal.force_retry('1s'::INTERVAL)
	Error:          pq: restart transaction: crdb_internal.force_retry(): TransactionRetryWithProtoRefreshError: forced by crdb_internal.force_retry()
	Time taken:     0.00076s
	Context:        context.Background

2021/12/08 21:39:58 upper/db: log_level=WARNING file=/home/anda/go/pkg/mod/github.com/upper/db/v4@v4.2.1/internal/sqladapter/session.go:646
	Session ID:     00007
	Transaction ID: 00006
	Query:          SELECT crdb_internal.force_retry('1s'::INTERVAL)
	Error:          pq: restart transaction: crdb_internal.force_retry(): TransactionRetryWithProtoRefreshError: forced by crdb_internal.force_retry()
	Time taken:     0.00087s
	Context:        context.Background

2021/12/08 21:39:58 Balances:
	accounts[717375400193490945]: 999
	accounts[717375400493318145]: 887
	accounts[717375400511373313]: 342
```
