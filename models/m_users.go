package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Users struct {
	Id       uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	Username string    `json:"username" gorm:"type:varchar(60);not null" firestore:"name"`
	Name     string    `json:"name" gorm:"type:varchar(60);not null" firestore:"name"`
	PhoneNo  string    `json:"phone_no" gorm:"type:varchar(15);Index:idx_phone,unique"`
	Email    string    `json:"email" gorm:"type:varchar(60);Index:idx_email,unique"`
	IsActive bool      `json:"is_active" gorm:"type:boolean"`
	JoinDate time.Time `json:"join_date" gorm:"type:timestamp with time zone;default:now()"`
	Password string    `json:"-" gorm:"type:varchar(200)"`
	Avatar   string    `json:"avatar" gorm:"varchar(100)" `
	Model
}
