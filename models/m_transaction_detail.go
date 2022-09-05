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
	TransactionId uuid.UUID `json:"transaction_id" gorm:"type:uuid;not null"`
	UserAppId     uuid.UUID `json:"user_app_id" gorm:"type:uuid;not null"`
	IsParent      bool      `json:"is_parent" gorm:"type:boolean;default:false"`
	ProductId     uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	ExtraAdult    int64     `json:"extra_adult" gorm:"type:integer"`
	Duration      int64     `json:"duration" gorm:"type:integer"`
	CheckIn       time.Time `json:"check_in" gorm:"type:timestamp"`
	CheckOut      time.Time `json:"check_out" gorm:"type:timestamp"`
	Amount        float64   `json:"amount" valid:"Required" gorm:"type:numeric(20,2)"`
}
