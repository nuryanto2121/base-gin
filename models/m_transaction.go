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
	CustomerId        uuid.UUID         `json:"customer_id" gorm:"type:uuid;not null"`
	OutletId          uuid.UUID         `json:"outlet_id" gorm:"type:uuid;not null"`
	TotalAmount       float64           `json:"total_amount"  gorm:"type:numeric(20,2)"`
	TotalTicket       int64             `json:"total_ticket"  gorm:"type:integer"`
	StatusPayment     StatusPayment     `json:"status_payment" gorm:"type:integer;not null"`
	StatusTransaction StatusTransaction `json:"status_transaction" gorm:"type:integer;not null"`
	PaymentId         uuid.UUID         `json:"payment_id" gorm:"type:uuid"` //update from online payment
	PaymentCode       PaymentCode       `json:"payment_code" gorm:"type:integer"`
	PaymentToken      uuid.UUID         `json:"payment_token" gorm:"type:uuid"`
	PaymentStatusDesc string            `json:"payment_status_desc" gorm:"type:varchar(250)"`
	Description       string            `json:"description"  gorm:"type:varchar(250)"`
}
