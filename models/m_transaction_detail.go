package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TransactionDetail struct {
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	AddTransactionDetail
	Model
}

type AddTransactionDetail struct {
	TransactionId        uuid.UUID         `json:"transaction_id" gorm:"type:uuid;not null"`
	TicketNo             string            `json:"ticket_no" gorm:"type:varchar(60)"`
	CustomerId           uuid.UUID         `json:"customer_id" gorm:"type:uuid;not null"`
	IsChildren           bool              `json:"is_children" gorm:"type:boolean;default:false"`
	ProductId            uuid.UUID         `json:"product_id" gorm:"type:uuid;not null"`
	ProductQty           int64             `json:"product_qty" gorm:"type:integer"`
	Duration             int64             `json:"duration" gorm:"type:integer"`
	CheckIn              time.Time         `json:"check_in" gorm:"type:timestamp"`
	CheckOut             time.Time         `json:"check_out" gorm:"type:timestamp"`
	Amount               float64           `json:"amount" valid:"Required" gorm:"type:numeric(20,2)"`
	Price                float64           `json:"price"  gorm:"type:numeric(20,2)"`
	FlagNotifSend        int64             `json:"flag_notif_send" gorm:"type:integer;default:0"`
	StatusTransactionDtl StatusTransaction `json:"status_transaction_dtl" gorm:"type:integer"`
}

type TransactionDetailRaw struct {
	TransactionDetail
	SkuName string `json:"sku_name"`
}
