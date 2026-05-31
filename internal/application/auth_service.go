package application

import (
	"context"
	"errors"
	"fmt"

	"alvintanoto.id/go-template/internal/domain/user"
	"go.uber.org/zap"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(plainPassword, encodedHash string) (bool, error)
}

type AuthService struct {
	log      *zap.Logger
	hasher   PasswordHasher
	userRepo user.Repository
}

func NewAuthService(log *zap.Logger, userRepo user.Repository, hasher PasswordHasher) *AuthService {
	return &AuthService{
		log:      log,
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

func (s *AuthService) ExecLogin(ctx context.Context, email, plainPassword string) (string, error) {
	domainUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.log.Error("find email", zap.Error(err))
		return "", errors.New("invalid email or password")
	}

	fmt.Println(plainPassword, domainUser.Password)
	match, err := s.hasher.Verify(plainPassword, domainUser.Password)
	if err != nil || !match {
		s.log.Error("verify", zap.Error(err))
		return "", errors.New("invalid email or password")
	}

	return "token", nil
}
