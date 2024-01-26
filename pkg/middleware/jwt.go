package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"go.uber.org/zap"
)

type JWTValidatorConfig struct {
	Audience string
	Domain   string
}

type CustomClaims struct {
	Permissions []string `json:"permissions"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func (c CustomClaims) HasPermission(expectedPermission string) bool {
	for _, permission := range c.Permissions {
		if permission == expectedPermission {
			return true
		}
	}

	return false
}

func EnsureValidToken(config JWTValidatorConfig) func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse(config.Domain)
	if err != nil {
		zap.L().Fatal("Failed to parse the issuer url: %v", zap.Error(err))
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{config.Audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)

	if err != nil {
		zap.L().Fatal("Failed to set up the jwt validator", zap.Error(err))
	}

	errorHandler := func(res http.ResponseWriter, req *http.Request, err error) {
		zap.L().Error("An error occurred during jwt validation", zap.Error(err))

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)
		_, _ = res.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

func RequirePermission(requiredPermission string) func(next http.Handler) http.Handler {
	// This returned func is the actual middleware created per route by invoking RequirePermission
	return func(next http.Handler) http.Handler {

		// This returned func is the handler invoked per request by the middleware
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := req.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			claims := token.CustomClaims.(*CustomClaims)

			if !claims.HasPermission(requiredPermission) {
				zap.L().Info(fmt.Sprintf("Subject %s is missing a required permission: %s", token.RegisteredClaims.Subject, requiredPermission))

				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"message":"Insufficient permissions."}`))

				return
			}

			next.ServeHTTP(w, req)
		})
	}
}
