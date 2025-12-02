package passutil

import (
	"github.com/apotourlyan/ludus-studii/pkg/errorutil"
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

type bcryptHasher struct{}

func NewHasher() Hasher {
	return &bcryptHasher{}
}

func (b *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), errorutil.SystemError(err)
}

func (b *bcryptHasher) Verify(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
