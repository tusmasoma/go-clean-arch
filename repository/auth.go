//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import "context"

type AuthRepository interface {
	GenerateToken(ctx context.Context, userID, email string) (jwt string, jti string)
	ValidateAccessToken(ctx context.Context, jwt string) (map[string]string, error)
	// GetPayloadFromToken(ctx context.Context, jwt string) (map[string]string, error)
}
