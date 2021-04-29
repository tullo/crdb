package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// Migrations
//
//go:embed db/migrations
var migrations embed.FS

func walkMigrations(log *log.Logger) error {
	return fs.WalkDir(migrations, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		log.Printf("path=%q, isDir=%v\n", path, d.IsDir())
		return nil
	})

}

// Remove container and linked volumes from docker.
func removeContainer(pool *dockertest.Pool, container *dockertest.Resource) {
	if err := pool.Purge(container); err != nil {
		log.Println("Could not purge container:", err)
	}
}

func startContainer(pool *dockertest.Pool) (*dockertest.Resource, error) {
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
		return nil, fmt.Errorf("could not start container: %w", err)
	}

	return container, nil
}

func connect(pool *dockertest.Pool, dsn string) (*sql.DB, error) {
	var db *sql.DB
	// Connect using exponential backoff-retry.
	if err := pool.Retry(func() error {
		log.Println("ping", dsn)
		var err error
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return db, nil
}

func main() {
	var (
		db  *sql.DB
		log = log.New(os.Stdout, "", log.Lshortfile)
	)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	container, err := startContainer(pool)
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("postgresql://root@0.0.0.0:%s/someshop?sslmode=disable", container.GetPort("26257/tcp"))
	db, err = connect(pool, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("DB close error:", err)
		}
	}()

	if _, err = db.Exec(`CREATE DATABASE IF NOT EXISTS someshop;
		CREATE USER IF NOT EXISTS johndoe; 
		GRANT ALL ON DATABASE someshop TO johndoe;`); err != nil {
		log.Println(err)
	}

	driver, err := cockroachdb.WithInstance(db, &cockroachdb.Config{})
	if err != nil {
		log.Println(err)
		removeContainer(pool, container)
		return
	}

	src, err := httpfs.New(http.FS(migrations), "db/migrations")
	if err != nil {
		log.Println(err)
		walkMigrations(log)
		removeContainer(pool, container)
		return
	}

	mig, err := migrate.NewWithInstance("httpfs", src, "someshop", driver)
	if err != nil {
		log.Println(err)
		removeContainer(pool, container)
		return
	}

	err = mig.Up()
	if err != nil {
		log.Println(err)
		removeContainer(pool, container)
		return
	}

	removeContainer(pool, container)
}
