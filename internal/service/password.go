package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	nanoid "github.com/matoous/go-nanoid/v2"
)

type HashedPassword struct {
	Hash string
	Salt string
}

const (
	MinLength  = 8
	MaxLength  = 32
	saltLength = 8
)

var (
	ErrShortPassword = errors.New("password is too short")
	ErrLongPassword  = errors.New("password is too long")
)

func (s *Service) ValidatePassword(password string) error {
	if len(password) < MinLength {
		return ErrShortPassword
	}

	if len(password) > MaxLength {
		return ErrLongPassword
	}

	return nil
}

func (s *Service) MakePassword(password string) (HashedPassword, error) {
	salt, err := nanoid.New(saltLength)
	if err != nil {
		return HashedPassword{}, err
	}

	return HashedPassword{
		Hash: hex.EncodeToString(makeHash(password, salt)),
		Salt: salt,
	}, nil
}

func (s *Service) VerifyPassword(hash, salt, password string) bool {
	return bytes.Equal([]byte(hash), makeHash(password, salt))
}

func makeHash(password, salt string) []byte {
	hasher := sha256.New()

	hasher.Write([]byte(salt))
	hasher.Write([]byte(password))

	return hasher.Sum(nil)
}
