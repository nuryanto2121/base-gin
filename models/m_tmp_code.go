package models

type TmpCode struct {
	ID     int64  `json:"id" gorm:"primary_key;auto_increment:true"`
	Prefix string `json:"prefix" gorm:"prefix"`
	SeqNo  int64  `json:"seq_no" gorm:"seq_no"`
	Code   string `json:"code" gorm:"code"`
	Model
}
