package models

import uuid "github.com/satori/go.uuid"

type SkuManagement struct {
	Id           uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	SkuName      string    `json:"sku_name" gorm:"type:varchar(60);Index:idx_skuname,unique;not null"`
	Duration     int64     `json:"duration" gorm:"type:integer"`
	PriceWeekDay float64   `json:"price_weekday" gorm:"type:numeric(20,2)"`
	PriceWeekEnd float64   `json:"price_weekend" gorm:"type:numeric(20,2)"`
	Model
}
