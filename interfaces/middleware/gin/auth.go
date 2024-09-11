package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/config"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/repository"
)

var ErrCacheMiss = errors.New("cache: key not found")

type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
}

type authMiddleware struct {
	ar repository.AuthRepository
}

func NewAuthMiddleware(ar repository.AuthRepository) AuthMiddleware {
	return &authMiddleware{
		ar: ar,
	}
}

func (am *authMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Info("Authentication failed: missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: missing Authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			log.Warn("Authorization failed: header format must be Bearer {token}")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed: header format must be Bearer {token}"})
			return
		}
		jwt := parts[1]

		if err := am.ar.ValidateAccessToken(jwt); err != nil {
			log.Warn("Authentication failed: invalid access token", log.Ferror(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Authentication failed: %v", err)})
			return
		}

		payload, err := am.ar.GetPayloadFromToken(jwt)
		if err != nil {
			log.Warn("Authentication failed: invalid access token", log.Ferror(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Authentication failed: %v", err)})
			return
		}

		ctx = context.WithValue(ctx, config.ContextUserIDKey, payload["userId"])
		c.Request = c.Request.WithContext(ctx)

		log.Info("Successfully Authentication", log.Fstring("userID", payload["userId"]))
		c.Next()
	}
}
