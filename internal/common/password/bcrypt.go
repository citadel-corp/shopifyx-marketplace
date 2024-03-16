package password

import (
	"errors"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var (
	costStr = os.Getenv("BCRYPT_SALT")
)

func Hash(plaintextPassword string) (string, error) {
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		return "", err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func Matches(plaintextPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
