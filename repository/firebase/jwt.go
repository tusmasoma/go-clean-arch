package firebase

import (
	"context"

	"github.com/tusmasoma/go-clean-arch/repository"
)

type authRepository struct {
	client *Client
}

func NewAuthRepository(client *Client) repository.AuthRepository {
	return &authRepository{
		client: client,
	}
}

func (ar *authRepository) GenerateToken(ctx context.Context, userID, _ string) (jwt string, jti string) { //nolint: nonamedreturns // this is a common pattern for returning multiple values
	token, err := ar.client.cli.CustomToken(ctx, userID)
	if err != nil {
		return "", ""
	}
	return token, ""
}

func (ar *authRepository) ValidateAccessToken(ctx context.Context, jwt string) (map[string]string, error) {
	var emptyPayload map[string]string
	token, err := ar.client.cli.VerifyIDToken(ctx, jwt)
	if err != nil {
		return emptyPayload, err
	}
	return map[string]string{
		"userId": token.UID,
	}, nil
}
