package models

import uuid "github.com/satori/go.uuid"

type UserOutlets struct {
	UserId   uuid.UUID `json:"id" gorm:"type:uuid"`
	OutletId uuid.UUID `json:"outlet_id" gorm:"type:uuid;not null"`
	Model
}
