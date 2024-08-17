package mongodb

import (
	"context"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tusmasoma/go-clean-arch/config"
)

// TODO: Implement similar transaction handling for MongoDB using mongo.Session.

type Client struct {
	cli *mongo.Client
	db  string
}

func NewMongoDB(ctx context.Context) (*Client, error) {
	cfg, err := config.NewMongoDB(ctx)
	if err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}
	api := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetServerAPIOptions(api)

	// Set authentication options if user and password are provided.
	if cfg.User != "" && cfg.Password != "" {
		opts.SetAuth(options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		})
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB instance: %w", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB instance: %w", err)
	}

	return &Client{
		cli: client,
		db:  cfg.Database,
	}, nil
}
