package models

import uuid "github.com/satori/go.uuid"

type UserGroup struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddUserGroup
	Model
}

type AddUserGroup struct {
	UserId  uuid.UUID `json:"user_id" gorm:"type:uuid;not null" `
	GroupId uuid.UUID `json:"group_id" gorm:"type:uuid;not null" `
}
