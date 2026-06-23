package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Plaintext string
	Hash      []byte
}

func (p *Password) HashPassword(passwordPlaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordPlaintext), 12)
	if err != nil {
		return err
	}

	p.Plaintext = ""
	p.Hash = hash

	return nil
}

func (p *Password) Matches(passwordPlaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(passwordPlaintext))
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
