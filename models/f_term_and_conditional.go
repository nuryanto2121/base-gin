package models

type TermAndConditionalForm struct {
	Id          string `json:"id"`
	Description string `json:"description"  valid:"Required"`
}
