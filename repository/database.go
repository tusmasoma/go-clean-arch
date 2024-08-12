//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import "context"

type TransactionRepository interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
