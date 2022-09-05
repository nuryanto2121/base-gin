package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Transaction struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddTransaction
	Model
}

type AddTransaction struct {
	TransactionCode   string            `json:"transaction_code" gorm:"type:varchar(25);not null"`
	TransactionDate   time.Time         `json:"transaction_date" gorm:"type:timestamp;not null"`
	OutletId          uuid.UUID         `json:"outlet_id" gorm:"type:uuid;not null"`
	TotalAmount       float64           `json:"total_amount"  gorm:"type:numeric(20,2)"`
	StatusPayment     StatusPayment     `json:"status_payment" gorm:"type:integer;not null"`
	StatusTransaction StatusTransaction `json:"status_transaction" gorm:"type:integer;not null"`
}
