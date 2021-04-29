package roachdbtest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

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

func TestRoachPing(t *testing.T) {
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

	t.Logf("Stats %+v", db.Stats())
	_, err = db.Exec("SELECT 1")
	if err != nil {
		t.Errorf("Could not query: %w", err)
	}

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(container); err != nil {
		t.Errorf("Could not purge resource: %s", err)
	}
}
