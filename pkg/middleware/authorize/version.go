package authorize

import (
	"app/models"
	"app/pkg/db"
	"app/pkg/setting"
	"context"
)

func GetVersion(ctx context.Context, V *models.AppVersion) (result *models.AppVersion, err error) {
	dbConn := db.NewDBdelegate(setting.DatabaseSetting.Debug)
	dbConn.Init()

	conn := dbConn.Get(ctx)
	err = conn.Where("device_type = ? ", V.DeviceType).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil

}
