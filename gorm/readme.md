# GORM

- https://gorm.io/
- https://pkg.go.dev/gorm.io/gorm

## Driver

`go get gorm.io/gorm@latest`

## Transaction Retry Support

https://github.com/cockroachdb/cockroach-go/tree/master/crdb/crdbgorm

## Bank Example

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

`DATABASE_URL='postgresql://johndoe@localhost:26257/bank?sslmode=disable' ./gorm`

```console
2022/07/06 13:13:13 Creating 5 new accounts...
2022/07/06 13:13:13 Accounts created.
Balance at '2022-07-06 13:13:13.542826662 +0200 CEST m=+0.151700383':
a3869d47-eeda-4e4c-9564-7ddfc5ac0868 4159
c22bfbe2-3dcc-495c-b89e-410cef601168 7987
d205f242-dec4-481c-a5e6-ac83b1d763ee 2181
dc35d3f6-617c-417e-ae21-2b98dc62a90d 8181
f5eedc33-c059-4a71-b51c-72acfc3b3f39 1947
2022/07/06 13:13:13 Transferring 100 from account dc35d3f6-617c-417e-ae21-2b98dc62a90d to account d205f242-dec4-481c-a5e6-ac83b1d763ee...
2022/07/06 13:13:13 Funds transferred.
Balance at '2022-07-06 13:13:13.554112895 +0200 CEST m=+0.162986613':
a3869d47-eeda-4e4c-9564-7ddfc5ac0868 4159
c22bfbe2-3dcc-495c-b89e-410cef601168 7987
d205f242-dec4-481c-a5e6-ac83b1d763ee 2281
dc35d3f6-617c-417e-ae21-2b98dc62a90d 8081
f5eedc33-c059-4a71-b51c-72acfc3b3f39 1947
2022/07/06 13:13:13 Deleting accounts created...
2022/07/06 13:13:13 Accounts deleted.
```
