package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tusmasoma/go-clean-arch/repository"
	"github.com/tusmasoma/go-tech-dojo/config"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

var ErrCacheMiss = errors.New("cache: key not found")

type AuthMiddleware interface {
	Authenticate(nextFunc http.Handler) http.Handler
}

type authMiddleware struct {
	ar repository.AuthRepository
}

func NewAuthMiddleware(ar repository.AuthRepository) AuthMiddleware {
	return &authMiddleware{
		ar: ar,
	}
}

func (am *authMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Info("Authentication failed: missing Authorization header")
			http.Error(w, "Authentication failed: missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			log.Warn("Authorization failed: header format must be Bearer {token}")
			http.Error(w, "Authorization failed: header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}
		jwt := parts[1]

		if err := am.ar.ValidateAccessToken(jwt); err != nil {
			log.Warn("Authentication failed: invalid access token", log.Ferror(err))
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}

		payload, err := am.ar.GetPayloadFromToken(jwt)
		if err != nil {
			log.Warn("Authentication failed: invalid access token", log.Ferror(err))
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, config.ContextUserIDKey, payload["userId"])

		log.Info("Successfully Authentication", log.Fstring("userID", payload["userId"]))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
