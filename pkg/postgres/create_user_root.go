package postgres

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/369-engineer/369backend/account/models"
	util "gitlab.com/369-engineer/369backend/account/pkg/utils"
)

const (
	createUser = `

	`
)

func Create() (err error) {
	var ctx = context.Background()
	data := &models.Users{}
	//check data is exist
	query := Conn.WithContext(ctx).Where("email = 'root@gmail.com'").First(data)
	err = query.Error
	if err != nil {
		fmt.Println("error postgres.Create() getdataby", err)
	}

	if data.Id == uuid.Nil {
		pass, _ := util.Hash("Playtopia12345")
		data = &models.Users{
			Name:     "root",
			PhoneNo:  "098765432",
			Email:    "root@gmail.com",
			IsActive: true,
			JoinDate: time.Now().Unix(),
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
