package gorm

import (
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

var (
	db        *gorm.DB
	mysqlPort string
)

func TestMain(m *testing.M) {
	var closeMySQL func()
	var err error

	db, mysqlPort, closeMySQL, err = startMySQL()
	defer closeMySQL()
	if err != nil {
		log.Error("Failed to start MySQL: %v", err)
	}

	m.Run()
}

func startMySQL() (*gorm.DB, string, func(), error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not connect to Docker: %w", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not ping Docker: %w", err)
	}

	runOptions := &dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_USER=root",
			"MYSQL_ROOT_PASSWORD=goCleanArc",
			"MYSQL_DATABASE=goCleanArcTestDB",
		},
		Cmd: []string{
			"--character-set-server=utf8mb4",
			"--collation-server=utf8mb4_unicode_ci",
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
		return nil, "", nil, fmt.Errorf("could not start resource: %w", err)
	}

	port := resource.GetPort("3306/tcp")

	err = pool.Retry(func() error {
		dsn := fmt.Sprintf("root:goCleanArc@(localhost:%s)/goCleanArcTestDB?charset=utf8mb4&parseTime=True", port)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		sqlDB, err := db.DB() //nolint:govet // err shadowing
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	})
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not connect to MySQL: %w", err)
	}

	log.Info("Start MySQL container 🐳")

	return db, port, func() { closeMySQL(db, pool, resource) }, nil
}

func closeMySQL(db *gorm.DB, pool *dockertest.Pool, resource *dockertest.Resource) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("Failed to get raw database connection: %v", err)
	}
	if err = sqlDB.Close(); err != nil {
		log.Error("Failed to close MySQL connection: %v", err)
	}

	if err = pool.Purge(resource); err != nil {
		log.Error("Failed to purge MySQL container: %v", err)
	}

	log.Info("Closed MySQL container 🐳")
}

func ValidateErr(t *testing.T, err error, wantErr error) {
	if (err != nil) != (wantErr != nil) {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	} else if err != nil && wantErr != nil && err.Error() != wantErr.Error() {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	}
}
