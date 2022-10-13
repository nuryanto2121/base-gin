package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type ReportForm struct {
	StartDate time.Time `json:"start_date" form:"start_date" valid:"Required"`
	EndDate   time.Time `json:"end_date" form:"end_date" valid:"Required"`
	OutletId  string    `json:"outlet_id" form:"outlet_id"`
}

type ReportResponse struct {
	OutletName           string    `json:"outletName" gorm:"outlet_name"`
	OutletCity           string    `json:"outletCity" gorm:"outlet_city"`
	TransactionID        uuid.UUID `json:"transactionId" gorm:"transaction_id"`
	TicketNo             string    `json:"ticketNo" gorm:"ticket_no"`
	TransactionCode      string    `json:"transactionCode" gorm:"transaction_code"`
	TransactionDate      string    `json:"transactionDate" gorm:"transaction_date"`
	TotalBooked          int64     `json:"totalBooked" gorm:"total_booked"`
	PaymentMethod        string    `json:"paymentMethod" gorm:"payment_method"`
	StatusTransaction    string    `json:"statusTransaction" gorm:"status_transaction"`
	SkuName              string    `json:"skuName" gorm:"sku_name"`
	NumOfKidsBooked      int64     `json:"numOfKidsBooked" gorm:"num_of_kids_booked"`
	NumOfKidsCheckIn     int64     `json:"numOfKidsCheckIn" gorm:"num_of_kids_check_in"`
	DeltaBookedVsCheckIn int64     `json:"deltaBookedVsCheckIn" gorm:"delta_booked_vs_check_in"`
	Qty                  int64     `json:"qty" gorm:"qty"`
	Amount               float64   `json:"amount" gorm:"amount"`
	TotalAmount          float64   `json:"totalAmount" gorm:"total_amount"`
}
