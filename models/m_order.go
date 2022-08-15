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
	OrderID   string    `json:"sku_name" valid:"Required" gorm:"type:varchar(25);Index:idx_orderid,unique;not null"`
	OrderDate time.Time `json:"order_date" gorm:"type:timestamp;default:now()"`
	OutletId  uuid.UUID `json:"outlet_id" gorm:"type:uuid;not null"`
	ProductId uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Status    string    `json:"status" gorm:"type:varchar(10)";default:'aa'`
}
