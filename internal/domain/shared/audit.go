package shared

import "time"

type AuditData struct {
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt *time.Time
}

func NewAuditData(Creator string) AuditData {
	now := time.Now()
	return AuditData{
		CreatedAt: now,
		CreatedBy: Creator,
		UpdatedAt: now,
		UpdatedBy: Creator,
		DeletedAt: nil,
	}
}

func (a *AuditData) DataChanged(Updater string) {
	a.UpdatedAt = time.Now().UTC()
	a.UpdatedBy = Updater
}
