package models

type TermAndConditionalForm struct {
	Id                  string `json:"id"`
	TermAndCondition    string `json:"term_and_condition"  valid:"Required"`
	KebijakanAndPrivacy string `json:"kebijakan_and_privacy" valid:"Required"`
}
