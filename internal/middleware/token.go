package middleware

import (
	"ProjectGolang/internal/entity"
	jwtPkg "ProjectGolang/pkg/jwt"
	"ProjectGolang/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"strings"
)

var (
	ErrUnauthorized = response.NewError(http.StatusUnauthorized, "unauthorized, access token invalid or expired")
)

const (
	AccessTokenSecret = "JWT_ACCESS_TOKEN_SECRET"
)

type tokenMiddleware struct {
}

func newTokenMiddleware() *tokenMiddleware {
	return &tokenMiddleware{}
}

func (m *middleware) NewTokenMiddleware(ctx *fiber.Ctx) error {
	clientIP := ctx.IP()
	authHeader := ctx.Get("Authorization")

	m.log.WithFields(logrus.Fields{
		"path":      ctx.Path(),
		"method":    ctx.Method(),
		"client_ip": clientIP,
		"headers":   ctx.GetReqHeaders(),
	}).Info("Incoming request")

	if authHeader == "" {
		m.log.Error("No Authorization header present")
		return ErrUnauthorized
	}

	headerParts := strings.Split(authHeader, " ")
	m.log.WithFields(logrus.Fields{
		"auth_type": headerParts[0],
		"has_token": len(headerParts) > 1,
	}).Debug("Authorization header check")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		m.log.Error("Invalid Authorization format - must start with 'Bearer '")
		return ErrUnauthorized
	}

	userToken, err := jwtPkg.VerifyTokenHeader(ctx, AccessTokenSecret)
	if err != nil {
		m.log.WithError(err).Error("Token verification failed")
		return ErrUnauthorized
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		m.log.Error("Could not extract token claims")
		return ErrUnauthorized
	}

	m.log.WithFields(logrus.Fields{
		"claim_keys":   reflect.ValueOf(claims).MapKeys(),
		"exp":          claims["exp"],
		"id_exists":    claims["id"] != nil,
		"email_exists": claims["email"] != nil,
	}).Debug("Token claims")

	user := entity.UserLoginData{
		ID:       claims["id"].(string),
		Email:    claims["email"].(string),
		Username: claims["username"].(string),
	}
	ctx.Locals("user", user)

	m.log.Info("Authentication successful")
	return ctx.Next()
}
