package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

type RoleForm struct {
	Role     string `json:"role"  valid:"Required"`
	RoleName string `json:"role_name"  valid:"Required"`
}

type UserRoleDesc struct {
	UserId   uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Role     string         `json:"role"`
	RoleName string         `json:"role_name"`
	Outlets  datatypes.JSON `json:"outlets"`
}
