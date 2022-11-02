package models

import uuid "github.com/satori/go.uuid"

func (SkuManagement) TableName() string {
	return "sku_management"
}

type SkuManagement struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddSkuManagement
	Model
}

type AddSkuManagement struct {
	SkuName    string    `json:"sku_name" valid:"Required" gorm:"type:varchar(60);Index:idx_skuname,unique;not null"`
	IsBracelet bool      `json:"is_bracelet" gorm:"type:boolean"`
	IsFree     bool      `json:"is_free" gorm:"type:boolean"`
	StatusDay  StatusDay `json:"status_day" gorm:"type:integer;default:4000003"`
	Duration   int64     `json:"duration" gorm:"type:integer"`
	Qty        int64     `json:"-" gorm:"-"`
	Price      float64   `json:"price" valid:"Required" gorm:"type:numeric(20,2)"`
	// Price float64 `json:"price" valid:"Required" gorm:"type:numeric(20,2)"`
	// PriceWeekend float64 `json:"price_weekend" valid:"Required" gorm:"type:numeric(20,2)"`
}
