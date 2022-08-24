package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type AuditLogs struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddAuditLog
	Model
}

type AddAuditLog struct {
	AuditDate time.Time `json:"audit_date" gorm:"type:timestamp;default:now()"`
	Username  string    `json:"username" gorm:"type:varchar(60);not null" `
	OutletId  uuid.UUID `json:"outlet_id" gorm:"type:uuid;not null"`
	ProductId uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Qty       int64     `json:"qty" gorm:"type:integer;default:0"`
	QtyChange int64     `json:"qty_change" gorm:"type:integer;default:0"`
	QtyDelta  int64     `json:"qty_delta" gorm:"type:integer;default:0"`
	Source    string    `json:"source" gorm:"type:varchar(60);not null"`
}
