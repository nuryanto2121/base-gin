package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type OutletForm struct {
	OutletName     string             `json:"outlet_name" valid:"Required"`
	OutletCity     string             `json:"outlet_city" valid:"Required"`
	OvertimeAmount float64            `json:"overtime_amount"`
	ToleransiTime  int64              `json:"toleransi_time"`
	OutletDetail   []*AddOutletDetail `json:"outlet_detail"`
}

type OutletList struct {
	UserId         uuid.UUID `json:"user_id"`
	Role           string    `json:"role"`
	OutletId       uuid.UUID `json:"outlet_id"`
	ProductId      uuid.UUID `json:"product_id"`
	InventoryId    uuid.UUID `json:"inventory_id"`
	OutletName     string    `json:"outlet_name"`
	OutletCity     string    `json:"outlet_city"`
	SkuName        string    `json:"sku_name"`
	Qty            int64     `json:"qty"`
	IsFree         bool      `json:"is_free"`
	IsBracelet     bool      `json:"is_bracelet"`
	StatusDay      StatusDay `json:"status_day"`
	Duration       int64     `json:"duration"`
	OvertimeAmount float64   `json:"overtime_amount"`
	ToleransiTime  int64     `json:"toleransi_time"`
	Price          float64   `json:"price"`
	OutletPrice    float64   `json:"outlet_price"`
}

type OutletLookupList struct {
	OutletId    uuid.UUID `json:"outlet_id"`
	ProductId   uuid.UUID `json:"product_id"`
	InventoryId uuid.UUID `json:"inventory_id"`
	OutletName  string    `json:"outlet_name"`
	OutletCity  string    `json:"outlet_city"`
	SkuName     string    `json:"sku_name"`
	Qty         int64     `json:"qty"`
	Duration    int64     `json:"duration"`
	IsBracelet  bool      `json:"is_bracelet"`
	IsFree      bool      `json:"is_free"`
	StatusDay   StatusDay `json:"status_day"`
	Price       float64   `json:"price"`
	OutletPrice float64   `json:"outlet_price"`
}

type OutletPriceProductRequest struct {
	OutletId        string    `json:"outlet_id" form:"outlet_id"`
	TransactionDate time.Time `json:"transaction_date" form:"transaction_date"`
}

type OutletPriceProductResponse struct {
	ProductId  uuid.UUID `json:"product_id"`
	SkuName    string    `json:"sku_name"`
	IsBracelet bool      `json:"is_bracelet"`
	IsFree     bool      `json:"is_free"`
	StatusDay  StatusDay `json:"status_day"`
	Duration   int64     `json:"duration"`
	Price      float64   `json:"price"`
}
