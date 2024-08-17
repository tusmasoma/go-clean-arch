package mongodb

import (
	"context"
	"fmt"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client    *mongo.Client
	mongoPort string
)

func TestMain(m *testing.M) {
	var closeMongoDB func()
	var err error

	client, mongoPort, closeMongoDB, err = startMongoDB()
	defer closeMongoDB()
	if err != nil {
		log.Error("Failed to start MongoDB: %v", err)
	}

	m.Run()
}

func startMongoDB() (*mongo.Client, string, func(), error) {
	// pwd, err := os.Getwd()
	// if err != nil {
	// 	log.Error("Failed to get current directory: %v", err)
	// 	return nil, "", nil, err
	// }

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
		Repository: "mongo",
		Tag:        "latest",
	}

	resource, err := pool.RunWithOptions(runOptions,
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
			// hc.Mounts = []docker.HostMount{
			// 	{
			// 		Type:   "bind",
			// 		Source: pwd + "/test/seed.js",
			// 		Target: "/docker-entrypoint-initdb.d/seed.js",
			// 	},
			// }
		},
	)
	if err != nil {
		log.Error("Could not start resource: %s", err)
		return nil, "", nil, err
	}

	port := resource.GetPort("27017/tcp")

	err = pool.Retry(func() error {
		uri := fmt.Sprintf("mongodb://localhost:%s", port)
		clientOptions := options.Client().ApplyURI(uri)

		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			return err
		}
		return client.Ping(context.Background(), nil)
	})
	if err != nil {
		log.Error("Could not connect to MongoDB: %s", err)
		return nil, "", nil, err
	}

	log.Info("start MongoDB container🐳")

	return client, port, func() { closeMongoDB(client, pool, resource) }, nil
}

func closeMongoDB(client *mongo.Client, pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Error("Failed to close MongoDB connection: %v", err)
	}

	if err := pool.Purge(resource); err != nil {
		log.Error("Failed to purge MongoDB container: %v", err)
	}

	log.Info("close MongoDB container🐳")
}

func ValidateErr(t *testing.T, err error, wantErr error) {
	if (err != nil) != (wantErr != nil) {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	} else if err != nil && wantErr != nil && err.Error() != wantErr.Error() {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	}
}
