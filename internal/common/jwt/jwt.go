package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	key     = []byte(os.Getenv("JWT_SECRET"))
	baseURL = os.Getenv("BASE_URL")

	ErrUnknownClaims = errors.New("unknown claims type")
)

func Sign(ttl time.Duration, subject string) (string, error) {
	now := time.Now()
	expiry := now.Add(ttl)
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
			Issuer:    baseURL,
			Audience:  jwt.ClaimStrings{baseURL},
			Subject:   subject,
		},
	)
	return t.SignedString(key)
}

func VerifyAndGetSubject(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return "", err
	}

	// Checking token validity
	if !token.Valid {
		return "", fmt.Errorf("token invalid")
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims.Subject, nil
	} else {
		return "", ErrUnknownClaims
	}
}
