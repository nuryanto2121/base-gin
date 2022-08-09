package models

import uuid "github.com/satori/go.uuid"

type OutletForm struct {
	OutletName   string             `json:"outlet_name" valid:"Required"`
	OutletCity   string             `json:"outlet_city" valid:"Required"`
	OutletDetail []*AddOutletDetail `json:"outlet_detail"`
}

type OutletList struct {
	UserId             uuid.UUID `json:"user_id"`
	Role               string    `json:"role"`
	OutletId           uuid.UUID `json:"outlet_id"`
	ProductId          uuid.UUID `json:"product_id"`
	InventoryId        uuid.UUID `json:"inventory_id"`
	OutletName         string    `json:"outlet_name"`
	OutletCity         string    `json:"outlet_city"`
	SkuName            string    `json:"sku_name"`
	Qty                int64     `json:"qty"`
	PriceWeekDay       float64   `json:"price_weekday"`
	PriceWeekEnd       float64   `json:"price_weekend"`
	OutletPriceWeekDay float64   `json:"outlet_price_weekday"`
	OutletPriceWeekEnd float64   `json:"outlet_price_weekend"`
}
