package models

import uuid "github.com/satori/go.uuid"

type Roles struct {
	Id       uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	Role     string    `json:"role" gorm:"type:varchar(10);Index:idx_role,unique;not null"`
	RoleName string    `json:"role_name" gorm:"type:varchar(150)"`
	Model
}
