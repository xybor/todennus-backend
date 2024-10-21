package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(secret string) ([]byte, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), HashingCost)
	if err != nil {
		return nil, Wrap(ErrUnknown, err.Error())
	}

	return hashedSecret, nil
}

func ValidatePassword(hashedSecret, secret string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(secret))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrMismatchedPassword
		}

		return Wrap(ErrUnknown, err.Error())
	}

	return nil
}
