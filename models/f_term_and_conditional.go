package models

import uuid "github.com/satori/go.uuid"

type TermAndConditionalForm struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"  valid:"Required"`
}
