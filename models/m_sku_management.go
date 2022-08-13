package models

import uuid "github.com/satori/go.uuid"

type SkuManagement struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddSkuManagement
	Model
}

type AddSkuManagement struct {
	SkuName      string  `json:"sku_name" valid:"Required" gorm:"type:varchar(60);Index:idx_skuname,unique;not null"`
	Duration     int64   `json:"duration" gorm:"type:integer"`
	Qty          int64   `json:"qty" gorm:"-"`
	PriceWeekDay float64 `json:"price_week_day" valid:"Required" gorm:"type:numeric(20,2)"`
	PriceWeekEnd float64 `json:"price_week_end" valid:"Required" gorm:"type:numeric(20,2)"`
}
