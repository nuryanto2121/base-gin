package models

import "time"

type HolidayForm struct {
	HolidayDate time.Time `json:"holiday_date"  valid:"Required"`
	Description string    `json:"description"  valid:"Required"`
}
