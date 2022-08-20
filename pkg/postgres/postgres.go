package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"app/models"
	authorize "app/pkg/middleware/authorize"
	"app/pkg/setting"
	util "app/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Conn *gorm.DB

func Setup() {
	now := time.Now()
	var err error

	connectionstring := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Name,
		setting.DatabaseSetting.Port)
	fmt.Printf("%s", connectionstring)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Jakarta"
	Conn, err = gorm.Open(postgres.Open(connectionstring), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   setting.DatabaseSetting.TablePrefix,
			SingularTable: true,
		},
		PrepareStmt: true,
		Logger:      newLogger,
		// Logger: logger.Default.LogMode(logger.Info),
		// DryRun: true,
	})

	if err != nil {
		log.Printf("connection.setup err : %v", err)
		panic(err)
	}

	// Conn.SingularTable(true)
	Conn.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	Conn.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	// Conn.Callback().Delete().Replace("gorm:delete", deleteCallback)

	sqlDB, err := Conn.DB()
	if err != nil {
		log.Printf("connection.setup DB err : %v", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	// sqlDB.Sing
	// Conn.DB().SetMaxIdleConns(10)
	// Conn.DB().SetMaxOpenConns(100)

	go autoMigrate()

	timeSpent := time.Since(now)
	log.Printf("Config DatabaseSetting is ready in %v", timeSpent)
}

func autoMigrate() {
	// Add auto migrate bellow this line
	Trx, err := Conn.DB()
	if err != nil {
		log.Printf("connection.setup autoMigrate err : %v", err)
		// panic(err)
	}
	rest, err := Trx.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	`)

	log.Printf("%v", rest)
	if err != nil {
		log.Printf(" : %v", err)
		panic(err)
	}

	log.Println("STARTING AUTO MIGRATE ")
	err = Conn.AutoMigrate(
		// models.Users{},
		models.Users{},
		models.UserSession{},
		models.UserRole{},
		models.UserOutlets{},
		models.SkuManagement{},
		models.Roles{},
		models.RoleOutlet{},
		authorize.AppVersion{},
		models.Holidays{},
		models.Outlets{},
		models.OutletDetail{},
		models.Inventory{},
		models.TermAndConditional{},
		models.Order{},
		authorize.AppVersion{},
	)
	if err != nil {
		log.Printf("\nAutoMigrate : %#v", err)
		panic(err)
	}

	go func() {
		// init data first
		InitUser()
		InitVersion()
	}()

	log.Println("FINISHING AUTO MIGRATE ")
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(r *gorm.DB) {
	var ctx = context.Background()
	if r.Statement.Error == nil {
		TimeInput := r.Statement.Schema.LookUpField("created_at")
		TimeInput.Set(ctx, r.Statement.ReflectValue, util.GetTimeNow())

		TimeEdit := r.Statement.Schema.LookUpField("updated_at")
		TimeEdit.Set(ctx, r.Statement.ReflectValue, util.GetTimeNow())
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(r *gorm.DB) {
	if r.Statement.Changed() {
		r.Statement.SetColumn("updated_at", util.GetTimeNow())
	}

}

// addExtraSpaceIfExist:
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
