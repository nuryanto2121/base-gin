package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Model :
type Model struct {
	CreatedBy uuid.UUID `json:"created_by" gorm:"type:uuid" cql:"created_by"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp" cql:"created_at"`
	UpdatedBy uuid.UUID `json:"updated_by" gorm:"type:uuid" cql:"updated_by"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp" cql:"updated_at"`
}
