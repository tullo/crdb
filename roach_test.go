package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var pool *dockertest.Pool

func TestMain(m *testing.M) {
	var err error
	// Uses sensible defaults:
	// - on windows   (tcp/http)
	// - on linux/osx (socket)
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	os.Exit(m.Run())
}

func TestRoachIntegration(t *testing.T) {
	options := dockertest.RunOptions{
		Name:       "crdb",
		Repository: "cockroachdb/cockroach",
		Tag:        "v20.2.8",
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port("26257/tcp"): {{HostIP: "", HostPort: "26257"}},
			docker.Port("8080/tcp"):  {{HostIP: "", HostPort: "8080"}},
		},
		Cmd: []string{"start-single-node", "--insecure", "--listen-addr=0.0.0.0"},
	}
	hostConfig := func(hc *docker.HostConfig) {
		// Auto remove stopped container.
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	}
	container, err := pool.RunWithOptions(&options, hostConfig)
	if err != nil {
		t.Errorf("Could not start container: %s", err)
	}

	var (
		db  *sql.DB
		dsn = fmt.Sprintf("postgresql://admin@0.0.0.0:%s?sslmode=disable", container.GetPort("26257/tcp"))
	)
	// Connect using exponential backoff-retry.
	if err := pool.Retry(func() error {
		t.Log("ping", dsn)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Errorf("Could not connect to database: %s", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Log("close error:", err)
		}
	}()

	t.Logf("Stats %+v", db.Stats())
	_, err = db.Exec("SELECT 1")
	if err != nil {
		t.Error("Could not query:", err)
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance INT)"); err != nil {
		t.Error("Could not create table:", err)
	}

	if _, err := db.Exec(
		"INSERT INTO accounts (id, balance) VALUES (1, 1000), (2, 250)"); err != nil {
		t.Error("Could not insert row:", err)
	}

	printBalances(t, db, "")

	// Run a transfer in a transaction.
	if err = crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {
		return transferFunds(tx,
			1,   /* from acct# */
			2,   /* to acct# */
			100, /* amount */
		)
	}); err != nil {
		t.Error("Could not execute transaction:", err)
	}

	printBalances(t, db, "after tx")

	// Remove container and linked volumes from docker.
	if err := pool.Purge(container); err != nil {
		t.Errorf("Could not purge resource: %s", err)
	}
}

func transferFunds(tx *sql.Tx, from int, to int, amount int) error {
	// Read the balance.
	var fromBalance int
	if err := tx.QueryRow(
		"SELECT balance FROM accounts WHERE id = $1", from).Scan(&fromBalance); err != nil {
		return err
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// Perform the transfer.
	if _, err := tx.Exec(
		"UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from); err != nil {
		return err
	}
	if _, err := tx.Exec(
		"UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to); err != nil {
		return err
	}
	return nil
}

func printBalances(t *testing.T, db *sql.DB, s string) {
	rows, err := db.Query("SELECT id, balance FROM accounts")
	if err != nil {
		t.Error("Could not insert query accounts:", err)
	}
	defer rows.Close()
	t.Log("Balances:", s)
	for rows.Next() {
		var id, balance int
		if err := rows.Scan(&id, &balance); err != nil {
			t.Error(err)
		}
		t.Logf("%d %d\n", id, balance)
	}
}
