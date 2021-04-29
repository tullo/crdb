package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/dhui/dktest"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	"github.com/golang-migrate/migrate/v4/dktesting"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// 	dt "github.com/golang-migrate/migrate/v4/database/testing"
)

const defaultPort = 26257

var (
	opts = dktest.Options{Cmd: []string{"start-single-node", "--insecure", "--listen-addr=0.0.0.0"}, PortRequired: true, ReadyFunc: isReady}
	// Released versions: https://www.cockroachlabs.com/docs/releases/
	specs = []dktesting.ContainerSpec{
		{ImageName: "cockroachdb/cockroach:v20.2.8", Options: opts},
	}
)

func createDB(t *testing.T, c dktest.ContainerInfo) {
	ip, port, err := c.Port(defaultPort)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://root@%v:%v?sslmode=disable", ip, port))
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}()

	if _, err = db.Exec("CREATE DATABASE migtest"); err != nil {
		t.Fatal(err)
	}
}

func TestMigrate(t *testing.T) {
	dktesting.ParallelTest(t, specs, func(t *testing.T, ci dktest.ContainerInfo) {
		createDB(t, ci)

		ip, port, err := ci.Port(26257)
		if err != nil {
			t.Fatal(err)
		}

		addr := fmt.Sprintf("cockroach://root@%v:%v/migtest?sslmode=disable", ip, port)
		m, err := migrate.New("file://db/migrations", addr)
		if err != nil {
			t.Fatal(err)
		}
		//dt.TestMigrate(t, m)
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	})
}

func isReady(ctx context.Context, c dktest.ContainerInfo) bool {
	ip, port, err := c.Port(defaultPort)
	if err != nil {
		log.Println("port error:", err)
		return false
	}

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://root@%v:%v?sslmode=disable", ip, port))
	if err != nil {
		log.Println("open error:", err)
		return false
	}
	if err := db.PingContext(ctx); err != nil {
		log.Println("ping error:", err)
		return false
	}
	if err := db.Close(); err != nil {
		log.Println("close error:", err)
	}
	return true
}
