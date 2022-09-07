package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TransactionForm struct {
	TransactionDate time.Time                `json:"transaction_date" valid:"Required"`
	OutletId        uuid.UUID                `json:"outlet_id" valid:"Required"`
	Details         []*TransactionDetailForm `json:"details"`
}

type TransactionDetailForm struct {
	CustomerId uuid.UUID `json:"customer_id"`
	ProductId  uuid.UUID `json:"product_id"`
	ProductQty int64     `json:"product_qty"`
	Duration   int64     `json:"duration"`
	Amount     float64   `json:"amount"`
}

type TransactionList struct {
	Name              string            `json:"name" gorm:"name"`
	PhoneNo           string            `json:"phone_no" gorm:"phone_no"`
	IsParent          bool              `json:"is_parent" gorm:"is_parent"`
	CheckIn           time.Time         `json:"check_in" gorm:"check_in"`
	CheckOut          time.Time         `json:"check_out" gorm:"check_out"`
	Duration          int64             `json:"duration" gorm:"duration"`
	StatusTransaction StatusTransaction `json:"status_transaction" gorm:"status_transaction"`
	StatusPayment     StatusPayment     `json:"status_payment" gorm:"status_payment"`
}

type TransactionResponse struct {
	TransactionDate time.Time                    `json:"transaction_date" `
	OutletName      string                       `json:"outlet_name"`
	OutletCity      string                       `json:"outlet_city"`
	TotalTicket     int64                        `json:"total_ticket"`
	TotalAmount     float64                      `json:"total_amount"`
	Details         []*TransactionDetailResponse `json:"details"`
}

type TransactionDetailResponse struct {
	CustomerName string  `json:"customer_name"`
	Description  string  `json:"description"`
	ProductQty   int64   `json:"product_qty"`
	Duration     int64   `json:"duration"`
	Amount       float64 `json:"amount"`
}
