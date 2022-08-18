package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserSession struct {
	Id          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	UserId      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Token       string    `json:"token" gorm:"type:varchar(255);not null"`
	ExpiredDate time.Time `json:"expired_date" gorm:"type:timestamp(0) without time zone"`
	Model
}
