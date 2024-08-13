package redis

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest"
)

var (
	client    *redis.Client
	redisPort string
)

func TestMain(m *testing.M) {
	var closeRedis func()
	var err error

	client, redisPort, closeRedis, err = startRedis()
	defer closeRedis()
	if err != nil {
		log.Println(err)
	}

	m.Run()
}

func startRedis() (*redis.Client, string, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Printf("Could not construct pool: %s\n", err)
		return nil, "", nil, err
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Printf("Could not connect to Docker: %s", err)
		return nil, "", nil, err
	}

	redisOptions := &dockertest.RunOptions{
		Repository: "redis",
		Tag:        "5.0",
		Env: []string{
			"REDIS_PASSWORD=",
		},
	}

	redisResource, err := pool.RunWithOptions(redisOptions)
	if err != nil {
		log.Printf("Could not start Redis resource: %s", err)
		return nil, "", nil, err
	}

	redisPort = redisResource.GetPort("6379/tcp")

	err = pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%s", redisPort),
			Password: "",
			DB:       0,
		})
		cmd := client.Ping(context.Background())
		_, err = cmd.Result()
		return err
	})
	if err != nil {
		log.Printf("Could not connect to Redis container: %s", err)
		return nil, "", nil, err
	}

	log.Println("start Redis container🐳")

	return client, redisPort, func() { closeRedis(client, pool, redisResource) }, nil
}

func closeRedis(client *redis.Client, pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := client.Close(); err != nil {
		log.Fatalf("Failed to close redis: %s", err)
	}

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Failed to purge resource: %s", err)
	}

	log.Println("close Redis container🐳")
}

func ValidateErr(t *testing.T, err error, wantErr error) {
	if (err != nil) != (wantErr != nil) {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	} else if err != nil && wantErr != nil && err.Error() != wantErr.Error() {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	}
}
