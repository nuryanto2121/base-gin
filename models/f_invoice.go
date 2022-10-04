package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type InvoiceRequest struct {
	Duration        time.Duration `json:"duration"`
	TransactionCode string        `json:"transaction_code"`
	TransactionDate time.Time     `json:"transaction_date"`
	TotalAmount     float64       `json:"total_amount"`
	Payment         PaymentType   `json:"payment"`
	Customer        Customer      `json:"customer"`
	Items           []Item        `json:"items"`
	Callback        *Callback     `json:"callback"`
}

type Customer struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
type Item struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	MerchantName string  `json:"merchant"`
	Description  string  `json:"description"`
	Qty          int     `json:"qty"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
}

type Callback struct {
	SuccessRedirectURL string `json:"success_redirect_url"`
	FailureRedirectURL string `json:"failure_redirect_url"`
}

type InvoiceResponse struct {
	TransactionID string
	PaymentToken  string
	PaymentURL    string
}

type InvPatchStockRequest struct {
	OutletId  uuid.UUID `json:"outlet_id"`
	ProductId uuid.UUID `json:"product_id"`
	Qty       int64     `json:"qty"`
}
