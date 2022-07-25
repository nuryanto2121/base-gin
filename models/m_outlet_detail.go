package models

import uuid "github.com/satori/go.uuid"

type OutletDetail struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddOutletDetail
	Model
}

type AddOutletDetail struct {
	OutletId           uuid.UUID `json:"outlet_id" gorm:"type:uuid;not null"`
	ProductId          uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	OutletPriceWeekday float32   `json:"outlet_price_weekday" gorm:"type:numeric(20,2)"`
	OutletPriceWeekend float32   `json:"outlet_price_weekend" gorm:"type:numeric(20,2)"`
}
