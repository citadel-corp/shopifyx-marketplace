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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return key, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.RegisteredClaims); ok {
		return claims.Subject, nil
	} else {
		return "", ErrUnknownClaims
	}
}
