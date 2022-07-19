package models

import uuid "github.com/satori/go.uuid"

type Groups struct {
	Id          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	GroupName   string    `json:"group_name" gorm:"type:varchar(60);Index:idx_groupname,unique;not null"`
	Description string    `json:"decsription" gorm:"type:varchar(150)"`
	Model
}
