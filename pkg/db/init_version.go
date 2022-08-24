package db

import (
	"app/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

func InitVersion(conn *gorm.DB) (err error) {
	var ctx = context.Background()
	data := &models.AppVersion{}
	//check data is exist
	query := conn.WithContext(ctx).First(data)
	err = query.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("error postgres.InitVersion() getdataby", err)
	}

	if err == gorm.ErrRecordNotFound {
		data := &models.AppVersion{
			DeviceType: "web",
			Version:    100,
			MinVersion: 100,
		}

		query := conn.WithContext(ctx).Create(data)
		err = query.Error
		if err != nil {
			fmt.Println("error postgres.InitVersion() save ", err)
		}
	}

	return nil
}
