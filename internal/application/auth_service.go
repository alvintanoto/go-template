package application

import (
	"context"
	"errors"

	"alvintanoto.id/go-template/internal/domain/user"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type AuthService struct {
	hasher   PasswordHasher
	userRepo user.Repository
}

func NewAuthService(userRepo user.Repository, hasher PasswordHasher) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (s *AuthService) SignUp(ctx context.Context, email, plainPassword string) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	hashedPassword, err := s.hasher.Hash(plainPassword)
	if err != nil {
		return errors.New("failed to process security credentials")
	}

	newUser, err := user.NewUser(email, hashedPassword)
	if err != nil {
		return errors.New("failed to create user data")
	}

	return s.userRepo.Create(ctx, newUser)
}
