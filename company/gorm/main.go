package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tullo/crdb/company/gorm/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	const BANK_DNS = "postgresql://admin@0.0.0.0:26257/company_gorm?sslmode=disable"

	db := setupDB(BANK_DNS)

	router := httprouter.New()

	server := NewServer(db)
	server.RegisterRouter(router)

	log.Fatal(http.ListenAndServe(":6543", router))
}

func setupDB(addr string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(addr))
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Migrate the schema
	db.AutoMigrate(&model.Customer{}, &model.Order{}, &model.Product{})

	return db
}
