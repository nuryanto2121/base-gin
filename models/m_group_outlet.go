package models

import uuid "github.com/satori/go.uuid"

type RoleOutlet struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddRoleOutlet
	Model
}

type AddRoleOutlet struct {
	Role     string    `json:"role" gorm:"type:varchar(10);not null"`
	OutletId uuid.UUID `json:"outlet_id" gorm:"type:uuid;not null"`
	UserId   uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
}
