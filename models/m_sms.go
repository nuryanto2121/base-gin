package models

import (
	uuid "github.com/satori/go.uuid"
)

type SmsLog struct {
	Id       uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	ToUserId uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	PhoneNo  string    `json:"phone_no" gorm:"type:varchar(20)"`
	Code     int64     `json:"code" gorm:"type:integer"`
	Message  string    `json:"token" gorm:"type:varchar(255);not null"`
	MoreInfo string    `json:"more_info" gorm:"type:varchar(150)"`
	Model
}
