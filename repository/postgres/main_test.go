package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	_ "github.com/lib/pq" // This blank import is used for its init function
)

var (
	db           *sql.DB
	postgresPort string
)

func TestMain(m *testing.M) {
	var closePostgres func()
	var err error

	db, postgresPort, closePostgres, err = startPostgres()
	defer closePostgres()
	if err != nil {
		log.Error("Failed to start PostgreSQL: %v", err)
	}

	m.Run()
}

func startPostgres() (*sql.DB, string, func(), error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Error("Failed to get current directory: %v", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Error("Could not connect to Docker: %s", err)
		return nil, "", nil, err
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Error("Could not ping Docker: %s", err)
		return nil, "", nil, err
	}

	runOptions := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_USER=root",
			"POSTGRES_PASSWORD=goCleanArc",
			"POSTGRES_DB=goCleanArcTestDB",
		},
	}

	resource, err := pool.RunWithOptions(runOptions,
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
			hc.Mounts = []docker.HostMount{
				{
					Type:   "bind",
					Source: pwd + "/init/postgresql.conf",
					Target: "/etc/postgresql/postgresql.conf",
				},
				{
					Type:   "bind",
					Source: pwd + "/test/dml.test.sql",
					Target: "/docker-entrypoint-initdb.d/dml.test.sql",
				},
				{
					Type:   "bind",
					Source: pwd + "/test/ddl.test.sql",
					Target: "/docker-entrypoint-initdb.d/ddl.test.sql",
				},
			}
		},
	)
	if err != nil {
		log.Error("Could not start resource: %s", err)
		return nil, "", nil, err
	}

	port := resource.GetPort("5432/tcp")

	err = pool.Retry(func() error {
		dsn := fmt.Sprintf("postgres://root:goCleanArc@localhost:%s/goCleanArcTestDB?sslmode=disable", port)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		log.Error("Could not connect to PostgreSQL: %s", err)
		return nil, "", nil, err
	}

	log.Info("start PostgreSQL container🐳")

	return db, port, func() { closePostgres(db, pool, resource) }, nil
}

func closePostgres(db *sql.DB, pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := db.Close(); err != nil {
		log.Error("Failed to close PostgreSQL connection: %v", err)
	}

	if err := pool.Purge(resource); err != nil {
		log.Error("Failed to purge PostgreSQL container: %v", err)
	}

	log.Info("close PostgreSQL container🐳")
}

func ValidateErr(t *testing.T, err error, wantErr error) {
	if (err != nil) != (wantErr != nil) {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	} else if err != nil && wantErr != nil && err.Error() != wantErr.Error() {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	}
}
