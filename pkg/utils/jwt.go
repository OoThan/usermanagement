package utils

import (
	"errors"
	"time"

	"github.com/OoThan/usermanagement/config"
	"github.com/OoThan/usermanagement/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Name string `json:"name"`
	Id   uint64 `json:"id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(name string, id uint64) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Name: name,
		Id:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(config.Rsa().PrivateKey)
	if err != nil {
		logger.Sugar.Error("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return config.Rsa().PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("access token is invalid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func GenerateRefreshToken(tokenString string) (string, error) {
	claims, err := ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
		return "", errors.New("token is not expired, yet")
	}

	refreshToken, err := GenerateAccessToken(claims.Name, claims.Id)
	if err != nil {
		return "", nil
	}

	return refreshToken, nil
}
