package models

import uuid "github.com/satori/go.uuid"

type UserGroup struct {
	Id        uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	UserId    uuid.UUID `json:"user_id" gorm:"type:type:uuid;not null" `
	GroupName string    `json:"group_name" gorm:"type:varchar(15);Index:idx_groupname,unique"`
	Model
}
