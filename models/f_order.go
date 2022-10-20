package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type OrderList struct {
	Id          uuid.UUID `json:"id"`
	OrderID     string    `json:"order_id"`
	OrderDate   time.Time `json:"order_date"`
	OutletName  string    `json:"outlet_name"`
	OrderLines  string    `json:"order_lines"`
	StartNumber int64     `json:"start_number"`
	EndNumber   int64     `json:"end_number"`
	Status      int64     `json:"status"`
}
