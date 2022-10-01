package models

import (
	uuid "github.com/satori/go.uuid"
)

type SmsLog struct {
	Id             uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	ToUserId       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	PhoneNo        string    `json:"phone_no" gorm:"type:varchar(20)"`
	Code           int64     `json:"code" gorm:"type:integer"`
	Message        string    `json:"token" gorm:"type:varchar(255);not null"`
	MessageTwillio string    `json:"message_twillio" gorm:"type:text"`
	Model
}

type TwillioResponseError struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	Status   int    `json:"status"`
}

type TwillioResponse struct {
	AccountSid          string          `json:"account_sid"`
	APIVersion          string          `json:"api_version"`
	Body                string          `json:"body"`
	DateCreated         string          `json:"date_created"`
	DateSent            string          `json:"date_sent"`
	DateUpdated         string          `json:"date_updated"`
	Direction           string          `json:"direction"`
	ErrorCode           interface{}     `json:"error_code"`
	ErrorMessage        interface{}     `json:"error_message"`
	From                string          `json:"from"`
	MessagingServiceSid string          `json:"messaging_service_sid"`
	NumMedia            string          `json:"num_media"`
	NumSegments         string          `json:"num_segments"`
	Price               interface{}     `json:"price"`
	PriceUnit           interface{}     `json:"price_unit"`
	Sid                 string          `json:"sid"`
	Status              string          `json:"status"`
	SubresourceUris     SubresourceUris `json:"subresource_uris"`
	To                  string          `json:"to"`
	URI                 string          `json:"uri"`
}

type SubresourceUris struct {
	Media string `json:"media"`
}
