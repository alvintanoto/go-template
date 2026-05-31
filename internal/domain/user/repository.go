package user

import (
	"context"

	"alvintanoto.id/go-template/internal/domain/shared"
)

type Repository interface {
	shared.CRUDRepository[User, string]

	ExistsByEmail(ctx context.Context, email string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
