package passutil

import (
	"strings"
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher_Hash_Success(t *testing.T) {
	cases := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "password123",
		},
		{
			name:     "complex password",
			password: "P@ssw0rd!#$%^&*()_+-=[]{}|;':,.<>?",
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 72),
		},
		{
			name:     "unicode password",
			password: "пароль密码🔐",
		},
		{
			name:     "empty string",
			password: "",
		},
	}

	hasher := NewHasher()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hash, err := hasher.Hash(c.password)
			testutil.DontWantError(t, err)
			testutil.DontWant(t, hash, "")
			testutil.DontWant(t, hash, c.password)
			testutil.WantPrefix(t, hash, "$2a$", "$2b$") // bcrypt prefix
		})
	}
}

func TestBcryptHasher_Hash_DifferentHashesForSamePassword(t *testing.T) {
	password := "samepassword"
	hasher := NewHasher()

	hash1, err1 := hasher.Hash(password)
	hash2, err2 := hasher.Hash(password)

	testutil.DontWantError(t, err1)
	testutil.DontWantError(t, err2)
	testutil.DontWant(t, hash1, hash2)
}

func TestBcryptHasher_Verify_Success(t *testing.T) {
	cases := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "password123",
		},
		{
			name:     "complex password",
			password: "P@ssw0rd!#$%^&*()_+-=[]{}|;':,.<>?",
		},
		{
			name:     "unicode password",
			password: "пароль密码🔐",
		},
		{
			name:     "empty string",
			password: "",
		},
	}

	hasher := NewHasher()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hash, err := hasher.Hash(c.password)
			testutil.DontWantError(t, err)

			err = hasher.Verify(c.password, hash)
			testutil.DontWantError(t, err)
		})
	}
}

func TestBcryptHasher_Verify_Failure(t *testing.T) {
	cases := []struct {
		name          string
		password      string
		wrongPassword string
	}{
		{
			name:          "completely different",
			password:      "correct",
			wrongPassword: "wrong",
		},
		{
			name:          "case sensitive",
			password:      "Password",
			wrongPassword: "password",
		},
		{
			name:          "extra character",
			password:      "password",
			wrongPassword: "password1",
		},
		{
			name:          "missing character",
			password:      "password",
			wrongPassword: "passwor",
		},
		{
			name:          "empty vs non-empty",
			password:      "password",
			wrongPassword: "",
		},
	}

	hasher := NewHasher()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hash, err := hasher.Hash(c.password)
			testutil.DontWantError(t, err)

			err = hasher.Verify(c.wrongPassword, hash)
			testutil.WantError(t, err)
			testutil.GotWant(t, err, bcrypt.ErrMismatchedHashAndPassword)
		})
	}
}

func TestBcryptHasher_Verify_InvalidHash(t *testing.T) {
	cases := []struct {
		name        string
		password    string
		invalidHash string
	}{
		{
			name:        "empty hash",
			password:    "password",
			invalidHash: "",
		},
		{
			name:        "invalid format",
			password:    "password",
			invalidHash: "not-a-valid-hash",
		},
		{
			name:        "truncated hash",
			password:    "password",
			invalidHash: "$2a$10$",
		},
		{
			name:        "plain text as hash",
			password:    "password",
			invalidHash: "password",
		},
	}

	hasher := NewHasher()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := hasher.Verify(c.password, c.invalidHash)
			testutil.WantError(t, err)
		})
	}
}

func TestBcryptHasher_DefaultCost(t *testing.T) {
	hasher := NewHasher()
	password := "testpassword"

	hash, err := hasher.Hash(password)
	testutil.DontWantError(t, err)

	cost, err := bcrypt.Cost([]byte(hash))
	testutil.DontWantError(t, err)
	testutil.GotWant(t, cost, bcrypt.DefaultCost)
}
