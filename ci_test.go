package main

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestWithGithubService(t *testing.T) {
	var (
		db  *sql.DB
		dsn = fmt.Sprintf("postgresql://postgres_user:postgres_password@%s:%d/postgres_db?sslmode=disable", "postgres", 5432)
	)
	// setup database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	} else if err := db.Ping(); err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		// query database
		firstQueryStart := time.Now()
		rows, err := db.Query("select 1;")
		firstQueryEnd := time.Now()
		if err != nil {
			panic(err)
		}

		// put the connection back to the pool so
		// that it can be reused by next iteration
		rows.Close()

		fmt.Println(fmt.Sprintf("query #%d took %s", i, firstQueryEnd.Sub(firstQueryStart).String()))
	}
}
