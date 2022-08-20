package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Order struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddOrder
	Model
}

type AddOrder struct {
	OrderID     string      `json:"sku_name" gorm:"type:varchar(25);Index:idx_orderid,unique;not null"`
	OrderDate   time.Time   `json:"order_date" valid:"Required" gorm:"type:timestamp;default:now()"`
	OutletId    uuid.UUID   `json:"outlet_id" valid:"Required" gorm:"type:uuid;not null"`
	ProductId   uuid.UUID   `json:"product_id" valid:"Required" gorm:"type:uuid;not null"`
	StartNumber int64       `json:"start_number" valid:"Required" gorm:"type:integer;not null"`
	EndNumber   int64       `json:"end_number" valid:"Required" gorm:"type:integer;not null"`
	Status      StatusOrder `json:"status" gorm:"type:integer;not null"`
}
