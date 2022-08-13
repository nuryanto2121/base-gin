package models

import uuid "github.com/satori/go.uuid"

type Inventory struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddInventory
	Model
}

type AddInventory struct {
	OutletId  uuid.UUID `json:"outlet_id" gorm:"type:uuid;not null"`
	ProductId uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Qty       int64     `json:"qty" gorm:"type:integer;default:0"`
	QtyChange int64     `json:"qty_change" gorm:"type:integer;default:0"`
	QtyDelta  int64     `json:"qty_delta" gorm:"type:integer;default:0"`
}
