package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Holidays struct {
	Id          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	HolidayDate time.Time `json:"holiday_date" gorm:"type:timestamp;default:now()"`
	Description string    `json:"description" gorm:"type:varchar(150)"`
	Model
}
