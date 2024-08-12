package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	_ "github.com/go-sql-driver/mysql" // This blank import is used for its init function
)

var (
	db        *sql.DB
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

// startMySQL はDockerを使用してMySQLコンテナを起動し、データベース接続を確立する関数です。
func startMySQL() (*sql.DB, string, func(), error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Error("Failed to get current directory: %v", err)
	}

	// Dockerのデフォルト接続方法を使用（Windowsではtcp/http、Linux/OSXではsocket）
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Error("Could not connect to Docker: %s", err)
		return nil, "", nil, err
	}

	// Dockerに接続を試みる
	err = pool.Client.Ping()
	if err != nil {
		log.Error("Could not ping Docker: %s", err)
		return nil, "", nil, err
	}

	// Dockerコンテナを起動する際に指定する設定定義
	runOptions := &dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_USERNAME=root",
			"MYSQL_ROOT_PASSWORD=goCleanArc",
			"MYSQL_DATABASE=goCleanArcTestDB",
		},
		Cmd: []string{
			"--character-set-server=utf8mb4",
			"--collation-server=utf8mb4_unicode_ci",
		},
	}

	// runOptions設定を適用してDockerコンテナを起動します。成功するとresourceは、起動したコンテナを表す。
	resource, err := pool.RunWithOptions(runOptions,
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
			hc.Mounts = []docker.HostMount{
				{
					Type:   "bind",
					Source: pwd + "/init/my.cnf",
					Target: "/etc/mysql/my.cnf",
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

	port := resource.GetPort("3306/tcp")

	// データベース接続が成功するまで定期的に接続試行を行うことを試みる(待機)
	err = pool.Retry(func() error {
		dsn := fmt.Sprintf("root:goCleanArc@(localhost:%s)/goCleanArcTestDB?charset=utf8mb4&parseTime=true", port)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		log.Error("Could not connect to MySQL: %s", err)
		return nil, "", nil, err
	}

	log.Info("start MySQL container🐳")

	// データベース接続とクリーンアップ関数を返却
	return db, port, func() { closeMySQL(db, pool, resource) }, nil
}

// closeMySQL はMySQLデータベースの接続を閉じ、Dockerコンテナを停止・削除する関数
func closeMySQL(db *sql.DB, pool *dockertest.Pool, resource *dockertest.Resource) {
	// データベース接続を切断
	if err := db.Close(); err != nil {
		log.Error("Failed to close MySQL connection: %v", err)
	}

	// Dockerコンテナを停止して削除
	if err := pool.Purge(resource); err != nil {
		log.Error("Failed to purge MySQL container: %v", err)
	}

	log.Info("close MySQL container🐳")
}

func ValidateErr(t *testing.T, err error, wantErr error) {
	if (err != nil) != (wantErr != nil) {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	} else if err != nil && wantErr != nil && err.Error() != wantErr.Error() {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	}
}
