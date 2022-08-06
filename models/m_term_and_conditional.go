package models

import uuid "github.com/satori/go.uuid"

type TermAndConditional struct {
	Id          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;DEFAULT:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"type:varchar(150)"`
	Description string    `json:"description" gorm:"type:varchar(150)"`
	Model
}
