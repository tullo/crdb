package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/dhui/dktest"
	"github.com/docker/docker/api/types/mount"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	"github.com/golang-migrate/migrate/v4/dktesting"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const defaultPort = 26257

var (
	opts = dktest.Options{
		Cmd: []string{
			"start-single-node",
			"--insecure",
			"--advertise-addr=0.0.0.0:26257",
		},
		PortRequired: true,
		LogStdout:    true,
		ReadyFunc:    isReady,
		Volumes:      []string{"crdb_data"},
		Mounts: []mount.Mount{
			mount.Mount{
				Source: "crdb_data",
				Target: "/cockroach/cockroach-data",
				Type:   "volume",
			},
		},
	}
	// Released versions: https://www.cockroachlabs.com/docs/releases/
	specs = []dktesting.ContainerSpec{
		{ImageName: "cockroachdb/cockroach:v22.1.2", Options: opts},
	}
)

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
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Println("ping error:", err)
		return false
	}

	return true
}

// Creates database migtest.
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

// Processes migration files located under db/migrations
// by applying all up migrations.
func TestMigrate(t *testing.T) {
	// Launch db container and call provided testFunc.
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
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	})
}
