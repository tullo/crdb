package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// Import GORM-related packages.
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Account is our model, which corresponds to the "accounts" database
// table.
type Account struct {
	ID      int `gorm:"primary_key"`
	Balance int
}

func transferFunds(db *gorm.DB, fromID int, toID int, amount int) error {
	var fromAccount Account
	var toAccount Account

	db.First(&fromAccount, fromID)
	db.First(&toAccount, toID)

	if fromAccount.Balance < amount {
		return fmt.Errorf("account %d balance %d is lower than transfer amount %d", fromAccount.ID, fromAccount.Balance, amount)
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount

	if err := db.Save(&fromAccount).Error; err != nil {
		return err
	}
	if err := db.Save(&toAccount).Error; err != nil {
		return err
	}
	return nil
}

func main() {
	db, err := gorm.Open(
		postgres.Open("postgresql://johndoe@0.0.0.0:26257/bank?sslmode=disable"),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Account{})

	// Insert two rows into the "accounts" table.
	var fromID = 1
	var toID = 2
	db.Create(&Account{ID: fromID, Balance: 1000})
	db.Create(&Account{ID: toID, Balance: 250})

	// Print balances before transfer.
	printBalances(db)

	var amount = 100

	// Transfer funds between accounts.  To handle potential
	// transaction retry errors, we wrap the call to `transferFunds`
	// in `crdbgorm.ExecuteTx`, a helper function for GORM which
	// implements a retry loop
	//func ExecuteTx(ctx context.Context, db *gorm.DB, opts *sql.TxOptions, fn func(tx *gorm.DB) error) error
	if err := crdbgorm.ExecuteTx(context.Background(), db, nil,
		func(tx *gorm.DB) error {
			return transferFunds(tx, fromID, toID, amount)
		},
	); err != nil {
		// For information and reference documentation, see:
		//   https://www.cockroachlabs.com/docs/stable/error-handling-and-troubleshooting.html
		fmt.Println(err)
	}

	printBalances(db)
	deleteAccounts(db)
}

func printBalances(db *gorm.DB) {
	var accounts []Account
	db.Find(&accounts)
	fmt.Printf("Balance at '%s':\n", time.Now())
	for _, account := range accounts {
		fmt.Printf("%d %d\n", account.ID, account.Balance)
	}
}

func deleteAccounts(db *gorm.DB) error {
	// Used to tear down the accounts table so we can re-run this
	// program.
	err := db.Exec("DELETE from accounts where ID > 0").Error
	if err != nil {
		return err
	}
	return nil
}
