package models

import uuid "github.com/satori/go.uuid"

type TermAndConditional struct {
	Id                  uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	TermAndCondition    string    `json:"term_and_condition" gorm:"type:text"`
	KebijakanAndPrivacy string    `json:"kebijakan_and_privacy" gorm:"type:text"`
	Model
}
