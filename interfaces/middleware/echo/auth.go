package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/config"

	"github.com/tusmasoma/go-clean-arch/repository"
)

type AuthMiddleware interface {
	Authenticate(next echo.HandlerFunc) echo.HandlerFunc
}

type authMiddleware struct {
	ar repository.AuthRepository
}

func NewAuthMiddleware(ar repository.AuthRepository) AuthMiddleware {
	return &authMiddleware{
		ar: ar,
	}
}

func (am *authMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			log.Info("Authentication failed: missing Authorization header")
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authentication failed: missing Authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			log.Warn("Authorization failed: header format must be Bearer {token}")
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authorization failed: header format must be Bearer {token}"})
		}
		jwt := parts[1]

		if err := am.ar.ValidateAccessToken(jwt); err != nil {
			log.Warn("Authentication failed: invalid access token", log.Ferror(err))
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": fmt.Sprintf("Authentication failed: %v", err)})
		}

		payload, err := am.ar.GetPayloadFromToken(jwt)
		if err != nil {
			log.Warn("Authentication failed: invalid access token", log.Ferror(err))
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": fmt.Sprintf("Authentication failed: %v", err)})
		}

		ctx = context.WithValue(ctx, config.ContextUserIDKey, payload["userId"])
		c.SetRequest(c.Request().WithContext(ctx))

		log.Info("Successfully Authentication", log.Fstring("userID", payload["userId"]))
		return next(c)
	}
}
