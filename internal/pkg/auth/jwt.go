package at

import (
	"errors"
	"fmt"
	"time"

	"GO-API/internal/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint   `json:"user_id"`
	ROle   string `json:"role"`
}

type JWTAuth struct {
	secretKey []byte
}

func NewJWTAuth(secretKey string) *JWTAuth {
	return &JWTAuth{
		secretKey: []byte(secretKey),
	}
}

func (j *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
	logger.Info("Validating token")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		logger.Error("Failed to parse token: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		logger.Error("Invalid token claims")
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil {
		if claims.ExpiresAt.Time.Before(time.Now()) {
			logger.Error("Token is expired")
			return nil, errors.New("token is expired")
		}
	}

	logger.Info("Token validated sccessfully for user: %d", claims.UserID)
	return claims, nil
}
