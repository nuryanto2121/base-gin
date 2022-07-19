package maria

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/369-engineer/369backend/account/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func Setup() {
	var (
		err error
		now = time.Now()
	)

	connectionstring := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Port,
		setting.DatabaseSetting.Name,
	)
	fmt.Printf("%s", connectionstring)
	Conn, err = gorm.Open(mysql.Open(connectionstring), &gorm.Config{})

	if err != nil {
		log.Printf("connection.setup err : %v", err)
		panic(err)
	}

	sqlDB, err := Conn.DB()
	if err != nil {
		log.Printf("connection.setup DB err : %v", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	go autoMigratePool()

	timeSpent := time.Since(now)
	log.Printf("Config database is ready in %v", timeSpent)

}

func autoMigratePool() {
	log.Println("STARTING AUTO MIGRATE ")
	err := Conn.AutoMigrate()
	if err != nil {
		log.Printf("Auoto Migrate err : %v", err)
		panic(err)
	}

	log.Println("FINISHING AUTO MIGRATE ")
}
