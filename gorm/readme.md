# Bank Example

[Covered steps](https://github.com/cockroachlabs/hello-world-go-gorm):
- Use  GORM to map Go-specific objects to SQL operations.
- Use `crdbgorm` package to handle transactions. Specifically:
  - `db.AutoMigrate(&Account{})` creates an `accounts` table based on the Account model.
  - `db.Create(&Account{})` inserts rows into the table.
  - `db.Find(&accounts)` selects from the table so that balances can be printed.
  - Transfer funds `transferFunds()`.
  - Handle retry errors `crdbgorm.ExecuteTx()`. 

Note that CockroachDB may require the client to retry a transaction in the case of read/write contention.

```sql
CREATE DATABASE bank;
CREATE USER johndoe WITH PASSWORD gopher;
GRANT ALL ON DATABASE bank TO johndoe;
-- DROP DATABASE bank CASCADE;
```

## GORM

```go
import (
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
```

`go build .`

`./gorm`

```console
Balance at '2021-04-28 11:10:11.868364981 +0200 CEST m=+0.058078433':
1 1000
2 250
Balance at '2021-04-28 11:10:11.877325459 +0200 CEST m=+0.067038908':
1 900
2 350
```
