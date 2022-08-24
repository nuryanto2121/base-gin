package models

type AppVersion struct {
	VersionID  int    `json:"version_id" gorm:"primary_key;auto_increment:true"`
	DeviceType string `json:"device_type" gorm:"type:varchar(20)" cql:"device_type"`
	Version    int    `json:"version" gorm:"type:integer" cql:"version"`
	MinVersion int    `json:"min_version" gorm:"type:integer" cql:"min_version"`
	Model
}
