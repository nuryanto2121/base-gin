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
	UserAppId  uuid.UUID `json:"user_app_id"`
	ExtraAdult int64     `json:"extra_adult"`
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
