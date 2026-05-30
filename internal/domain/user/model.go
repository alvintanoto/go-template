package user

import (
	"errors"
	"strings"

	"alvintanoto.id/go-template/internal/domain/shared"
	"github.com/google/uuid"
)

type User struct {
	UserID   string
	Email    string
	Password string

	shared.AuditData
}

func NewUser(email, passwordHash string) (*User, error) {
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email address")
	}

	// Generate UUID v7 directly at creation time
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("failed to generate user identity")
	}

	return &User{
		UserID:   id.String(),
		Email:    email,
		Password: passwordHash,
	}, nil
}
