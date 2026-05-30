package user

import (
	"context"

	"alvintanoto.id/go-template/internal/domain/shared"
)

type Repository interface {
	shared.CRUDRepository[User, int64]

	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
