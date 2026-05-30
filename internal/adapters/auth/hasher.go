package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"alvintanoto.id/go-template/internal/application"
	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 1
	saltLength  = 16
	keyLength   = 32
)

type Hasher struct {
}

func NewHasher() application.PasswordHasher {
	return &Hasher{}
}

func (b *Hasher) Hash(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedString := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash)

	return encodedString, nil
}
