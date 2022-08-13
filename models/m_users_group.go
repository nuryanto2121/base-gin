package models

import uuid "github.com/satori/go.uuid"

type UserRole struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddUserRole
	Model
}

type AddUserRole struct {
	UserId uuid.UUID `json:"user_id" gorm:"type:uuid;not null" `
	Role   string    `json:"role" gorm:"type:varchar(10);not null" `
}
