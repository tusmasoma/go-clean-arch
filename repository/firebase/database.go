package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"google.golang.org/api/option"

	"github.com/tusmasoma/go-clean-arch/config"
)

type Client struct {
	cli    *auth.Client
	apiKey string
}

func NewClient(ctx context.Context) (*Client, error) {
	conf, err := config.NewFirebaseConfig(ctx)
	if err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}

	opt := option.WithCredentialsFile(conf.CredentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Error("Failed to create firebase app", log.Ferror(err))
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Error("Failed to create firebase auth client", log.Ferror(err))
		return nil, err
	}
	return &Client{
		cli:    client,
		apiKey: conf.APIKey,
	}, nil
}

func (c *Client) postCall(ctx context.Context, url string, reqBody any, respBody any) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Error("Failed to marshal request body", log.Ferror(err))
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error("Failed to create request", log.Ferror(err))
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Failed to send request", log.Ferror(err))
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}
	return nil
}
