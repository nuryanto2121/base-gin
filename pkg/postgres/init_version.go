package postgres

import (
	authorize "app/pkg/middleware/authorize"
	"context"
	"fmt"

	"gorm.io/gorm"
)

func InitVersion() (err error) {
	var ctx = context.Background()
	data := &authorize.AppVersion{}
	//check data is exist
	query := Conn.WithContext(ctx).First(data)
	err = query.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("error postgres.InitVersion() getdataby", err)
	}

	if err == gorm.ErrRecordNotFound {
		data := &authorize.AppVersion{
			DeviceType: "web",
			Version:    100,
			MinVersion: 100,
		}

		query := Conn.WithContext(ctx).Create(data)
		err = query.Error
		if err != nil {
			fmt.Println("error postgres.InitVersion() save ", err)
		}
	}

	return nil
}
