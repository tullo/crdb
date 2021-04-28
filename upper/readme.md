# Bank Example

[Covered steps](https://github.com/cockroachlabs/hello-world-go-upperdb):
- Use  `upper/db` to map Go-specific objects to SQL operations.
- Creates the accounts table, if it does not already exist.
- Deletes any existing rows in the accounts table.
- Inserts two rows into the accounts table.
- Prints the rows in the accounts table to the terminal.
- Deletes the first row in the accounts table.
- Updates the rows in the accounts table within an explicit [transaction](https://www.cockroachlabs.com/docs/v20.2/transactions).
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

`./upper`

```console
2021/04/28 12:16:26 Balances:
	accounts[653834932749434881]: 1000
	accounts[653834932766900225]: 250
2021/04/28 12:16:26 Balances:
	accounts[653834932749434881]: 500
	accounts[653834932766900225]: 999
2021/04/28 12:16:26 upper/db: log_level=WARNING file=/home/anda/go/pkg/mod/github.com/upper/db/v4@v4.1.0/internal/sqladapter/session.go:643
	Session ID:     00006
	Transaction ID: 00005
	Query:          SELECT crdb_internal.force_retry('1s'::INTERVAL)
	Error:          pq: restart transaction: crdb_internal.force_retry(): TransactionRetryWithProtoRefreshError: forced by crdb_internal.force_retry()
	Time taken:     0.00029s
	Context:        context.Background

2021/04/28 12:16:26 upper/db: log_level=WARNING file=/home/anda/go/pkg/mod/github.com/upper/db/v4@v4.1.0/internal/sqladapter/session.go:643
	Session ID:     00007
	Transaction ID: 00006
	Query:          SELECT crdb_internal.force_retry('1s'::INTERVAL)
	Error:          pq: restart transaction: crdb_internal.force_retry(): TransactionRetryWithProtoRefreshError: forced by crdb_internal.force_retry()
	Time taken:     0.00026s
	Context:        context.Background

2021/04/28 12:16:26 upper/db: log_level=WARNING file=/home/anda/go/pkg/mod/github.com/upper/db/v4@v4.1.0/internal/sqladapter/session.go:643
	Session ID:     00008
	Transaction ID: 00007
	Query:          SELECT crdb_internal.force_retry('1s'::INTERVAL)
	Error:          pq: restart transaction: crdb_internal.force_retry(): TransactionRetryWithProtoRefreshError: forced by crdb_internal.force_retry()
	Time taken:     0.00029s
	Context:        context.Background

2021/04/28 12:16:26 Balances:
	accounts[653834932766900225]: 999
	accounts[653834933128921089]: 887
	accounts[653834933139439617]: 342
```
