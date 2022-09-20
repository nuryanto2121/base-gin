package models

import (
	"time"
)

type Child struct {
	ChildrenId string    `json:"children_id"`
	Name       string    `json:"name"`
	DOB        time.Time `json:"dob"`
}

type ChildForm struct {
	Childs []*Child `json:"child" validate:"dive"`
}
