package models

import uuid "github.com/satori/go.uuid"

type TermAndConditionalForm struct {
	Id          uuid.UUID `json:"term_and_conditional_id" gorm:"type:uuid;not null"`
	Description string    `json:"description"  valid:"Required"`
}
