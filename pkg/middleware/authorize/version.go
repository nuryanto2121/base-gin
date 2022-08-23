package authorize

import (
	"app/models"
	"context"

	"gorm.io/gorm"
)

type AppVersion struct {
	VersionID  int    `json:"version_id" gorm:"primary_key;auto_increment:true"`
	DeviceType string `json:"device_type" gorm:"type:varchar(20)" cql:"device_type"`
	Version    int    `json:"version" gorm:"type:integer" cql:"version"`
	MinVersion int    `json:"min_version" gorm:"type:integer" cql:"min_version"`
	models.Model
}

func (V *AppVersion) GetVersion(ctx context.Context, Conn *gorm.DB) (result *AppVersion, err error) {
	err = Conn.Where("device_type = ? ", V.DeviceType).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil

}
