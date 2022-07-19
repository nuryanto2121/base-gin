package postgres

import (
	"context"
	"fmt"
	"time"

	"app/models"
	util "app/pkg/utils"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func Create() (err error) {
	var ctx = context.Background()
	data := &models.Users{}
	//check data is exist
	query := Conn.WithContext(ctx).Where("email = 'root@gmail.com'").First(data)
	err = query.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("error postgres.Create() getdataby", err)
	}

	if data.Id == uuid.Nil {
		pass, _ := util.Hash("Playtopia12345")
		data = &models.Users{
			Username: "root",
			Name:     "root",
			PhoneNo:  "098765432",
			Email:    "root@gmail.com",
			IsActive: true,
			JoinDate: time.Now(),
			Password: pass,
		}

		query := Conn.WithContext(ctx).Create(data)
		err = query.Error
		if err != nil {
			fmt.Println("error postgres.Create() save ", err)
		}
	}

	return nil
}
