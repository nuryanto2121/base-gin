package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

type GroupForm struct {
	GroupCode   string `json:"GroupCode"  valid:"Required"`
	Description string `json:"description"  valid:"Required"`
}

type UserGroupDesc struct {
	GroupId     uuid.UUID      `json:"group_id" gorm:"type:uuid;not null"`
	GroupCode   string         `json:"GroupCode"`
	Description string         `json:"description"`
	Outlets     datatypes.JSON `json:"outlets"`
}
