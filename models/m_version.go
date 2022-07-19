package models

type VersionApps struct {
	VersionID int    `json:"version_id" gorm:"PRIMARY_KEY"`
	OS        string `json:"os" gorm:"type:varchar(20)"`
	Version   int    `json:"version" gorm:"type:integer"`
}
