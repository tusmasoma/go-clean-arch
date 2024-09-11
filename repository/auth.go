//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

type AuthRepository interface {
	GenerateToken(userID, email string) (jwt string, jti string)
	ValidateAccessToken(jwt string) error
	GetPayloadFromToken(jwt string) (map[string]string, error)
}
