package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tullo/crdb/company/gobun/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func main() {
	db := setupDB()
	defer db.Close()

	server := NewServer(db)
	router := httprouter.New()
	server.RegisterRouter(router)

	log.Fatal(http.ListenAndServe(":6543", router))
}

func setupDB() *bun.DB {
	// pgconn := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr("0.0.0.0:26257"),
		pgdriver.WithUser("admin"),
		pgdriver.WithPassword(""),
		pgdriver.WithDatabase("company_bun"),
		pgdriver.WithApplicationName("testcompany"),
		pgdriver.WithInsecure(true),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
	)

	db := bun.NewDB(sql.OpenDB(pgconn), pgdialect.New())

	db.RegisterModel((*model.OrderToProduct)(nil))

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	models := []interface{}{
		(*model.Customer)(nil),
		(*model.Order)(nil),
		(*model.Product)(nil),
		(*model.OrderToProduct)(nil),
	}
	if err := db.ResetModel(context.Background(), models...); err != nil {
		panic(fmt.Sprintf("failed to create a table: %v", err))
	}

	return db
}
