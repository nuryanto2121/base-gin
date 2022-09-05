package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserApps struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddUserApps
	Model
}

type AddUserApps struct {
	Name     string    `json:"name" gorm:"type:varchar(60);not null" firestore:"name"`
	PhoneNo  string    `json:"phone_no" gorm:"type:varchar(15)"` //;Index:idx_phone,unique"
	IsParent bool      `json:"is_parent" gorm:"type:boolean"`
	ParentId uuid.UUID `json:"parent_id" gorm:"type:uuid"`
	JoinDate time.Time `json:"join_date" gorm:"type:timestamp;default:now()"`
	DOB      time.Time `json:"dob" gorm:"type:timestamp;default:now()"`
	Password string    `json:"-" gorm:"type:varchar(250)"`
	Avatar   string    `json:"avatar" gorm:"varchar(100)" `
}

//OutletId    uuid.UUID   `json:"outlet_id" valid:"Required" gorm:"type:uuid;not null"`
