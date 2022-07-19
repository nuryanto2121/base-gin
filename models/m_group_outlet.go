package models

import uuid "github.com/satori/go.uuid"

type GroupOutlet struct {
	Id       uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	GroupId  uuid.UUID `json:"group_id" gorm:"type:uuid;not null" `
	OutletId uuid.UUID `json:"outlet_id" gorm:"type:uuid;not null" `
	Model
}
