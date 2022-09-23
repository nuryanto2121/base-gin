package models

import uuid "github.com/satori/go.uuid"

type MidtransResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	StatusCode    string   `json:"status_code,omitempty"`
	ErrorMessages []string `json:"error_messages,omitempty"`
}

type MidtransNotification struct {
	VaNumbers         []VaNumber      `json:"va_numbers"`
	TransactionTime   string          `json:"transaction_time"`
	TransactionStatus string          `json:"transaction_status"`
	TransactionID     string          `json:"transaction_id"`
	StatusMessage     string          `json:"status_message"`
	StatusCode        string          `json:"status_code"`
	SignatureKey      string          `json:"signature_key"`
	SettlementTime    string          `json:"settlement_time"`
	PaymentType       string          `json:"payment_type"`
	PaymentAmounts    []PaymentAmount `json:"payment_amounts"`
	OrderID           string          `json:"order_id"`
	MerchantID        string          `json:"merchant_id"`
	GrossAmount       string          `json:"gross_amount"`
	FraudStatus       string          `json:"fraud_status"`
	Currency          string          `json:"currency"`
	Acquirer          string          `json:"acquirer"`
	Issuer            string          `json:"issuer"`
	Store             string          `json:"store"`
}

type PaymentAmount struct {
	PaidAt string `json:"paid_at"`
	Amount string `json:"amount"`
}

type VaNumber struct {
	VaNumber string `json:"va_number"`
	Bank     string `json:"bank"`
}

type MidtransNotificationLog struct {
	Id                uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	VaNumbers         string    `json:"va_numbers"`
	TransactionTime   string    `json:"transaction_time"`
	TransactionStatus string    `json:"transaction_status"`
	TransactionID     string    `json:"transaction_id"`
	StatusMessage     string    `json:"status_message"`
	StatusCode        string    `json:"status_code"`
	SignatureKey      string    `json:"signature_key"`
	SettlementTime    string    `json:"settlement_time"`
	PaymentType       string    `json:"payment_type"`
	PaymentAmounts    string    `json:"payment_amounts"`
	OrderID           string    `json:"order_id"`
	MerchantID        string    `json:"merchant_id"`
	GrossAmount       string    `json:"gross_amount"`
	FraudStatus       string    `json:"fraud_status"`
	Currency          string    `json:"currency"`
	Acquirer          string    `json:"acquirer"`
	Issuer            string    `json:"issuer"`
	Store             string    `json:"store"`
}
