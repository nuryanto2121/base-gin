package models

type GroupForm struct {
	GroupCode   string `json:"GroupCode"  valid:"Required"`
	Description string `json:"description"  valid:"Required"`
}
