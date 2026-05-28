package user

import "alvintanoto.id/go-template/internal/domain/shared"

type User struct {
	UserID   string
	Email    string
	Password string

	shared.AuditData
}
