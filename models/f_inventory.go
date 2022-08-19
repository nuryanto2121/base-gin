package models

import uuid "github.com/satori/go.uuid"

type InventoryForm struct {
	OutletId  uuid.UUID `json:"outlet_id" valid:"Required"`
	ProductId uuid.UUID `json:"product_id" valid:"Required"`
	Qty       int64     `json:"qty"`
	QtyChange int64     `json:"qty_change" valid:"Required"`
}