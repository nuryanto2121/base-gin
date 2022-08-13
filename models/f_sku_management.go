package models

type SkuMgmForm struct {
	SkuName      string  `json:"sku_name" valid:"Required"`
	Duration     int64   `json:"duration"`
	PriceWeekday float64 `json:"price_weekday" valid:"Required"`
	PriceWeekend float64 `json:"price_weekend" valid:"Required"`
}
