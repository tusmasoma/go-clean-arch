package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tusmasoma/go-clean-arch/pkg/config"
	"github.com/tusmasoma/go-clean-arch/pkg/jwt"
)

var ErrCacheMiss = errors.New("cache: key not found")

type AuthMiddleware interface {
	Authenticate(nextFunc http.Handler) http.Handler
}

type authMiddleware struct {
	jwtGenerator jwt.Generator
}

func NewAuthMiddleware(jwtGenerator jwt.Generator) AuthMiddleware {
	return &authMiddleware{
		jwtGenerator: jwtGenerator,
	}
}

func (am *authMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authentication failed: missing Authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Authorization failed: header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}
		jwt := parts[1]
		if err := am.jwtGenerator.ValidateAccessToken(jwt); err != nil {
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}
		payload, err := am.jwtGenerator.GetPayloadFromToken(jwt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, config.ContextUserIDKey, payload["userId"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
