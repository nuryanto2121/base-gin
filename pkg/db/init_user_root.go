package db

import (
	"app/models"
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitUser(conn *gorm.DB) (err error) {
	var ctx = context.Background()
	data := &models.Users{}
	//check data is exist
	query := conn.WithContext(ctx).Where("email = 'root@gmail.com'").First(data)
	err = query.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("error postgres.InitUser() getdataby", err)
	}

	if data.Id == uuid.Nil {
		pass, _ := hash("Playtopia12345")
		// pass := ""
		data = &models.Users{
			Username: "root",
			Name:     "root",
			PhoneNo:  "098765432",
			Email:    "root@gmail.com",
			IsActive: true,
			JoinDate: time.Now(),
			Password: pass,
		}

		query := conn.WithContext(ctx).Create(data)
		err = query.Error
		if err != nil {
			fmt.Println("error postgres.InitUser() save ", err)
		}
	}

	return nil
}

// Hash :
func hash(text string) (string, error) {
	pwd := []byte(text)

	hashedPwd, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return text, err
	}
	return string(hashedPwd), nil
}
