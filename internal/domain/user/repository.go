package user

import (
	"alvintanoto.id/go-template/internal/domain/shared"
)

type Repository interface {
	shared.CRUDRepository[User, int64]
}
