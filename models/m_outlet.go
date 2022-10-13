package models

import uuid "github.com/satori/go.uuid"

type Outlets struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddOutlets
	Model
}

type AddOutlets struct {
	OutletName     string  `json:"outlet_name" gorm:"type:varchar(100);not null"`
	OutletCity     string  `json:"outlet_city" gorm:"type:varchar(60)"`
	OvertimeAmount float64 `json:"overtime_amount" gorm:"type:numeric(20,2)"`
	ToleransiTime  int64   `json:"toleransi_time" gorm:"type:integer"`
}
