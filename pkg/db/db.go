package db

import (
	"app/pkg/setting"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBGormDelegate interface {
	Init()
	InitNoUse()
	Get(ctx context.Context) *gorm.DB
	BeginTx() *gorm.DB
	Rollback()
	Commit()
}

type dbDelegate struct {
	dbGorm *gorm.DB
	once   sync.Once
	debug  bool
	tx     *gorm.DB
}

func NewDBdelegate(debug bool) DBGormDelegate {
	return &dbDelegate{
		debug: debug,
	}
}

// Init mysql client specific db
func (dbdget *dbDelegate) Init() {
	dbdget.run(true)
}

// InitNoUse client not specific db
func (dbdget *dbDelegate) InitNoUse() {
	dbdget.run(false)
}

func (dbdget *dbDelegate) run(withDB bool) {
	dbdget.once.Do(func() {
		var logLevel logger.LogLevel
		if setting.DatabaseSetting.Debug {
			logLevel = logger.Info
		} else {
			logLevel = logger.Silent
		}

		var err error
		var dbGorm *gorm.DB

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logLevel,    // Log level
				Colorful:      true,        // Disable color
			},
		)

		dbGorm, err = gorm.Open(postgres.Open(connectionstring()), &gorm.Config{
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
			//panic("init mysql failed: " + err.Error())
		}

		dbdget.dbGorm = dbGorm

		if dbdget.debug {
			dbdget.dbGorm = dbdget.dbGorm.Debug()
		}
	})
}

func (dbdget *dbDelegate) Get(ctx context.Context) *gorm.DB {
	tx := ctx.Value("tx")
	if tx != nil {
		return tx.(*gorm.DB)
	}

	return dbdget.dbGorm
}

// new transactions

func (dbdget *dbDelegate) BeginTx() *gorm.DB {
	return dbdget.dbGorm.Begin()
}

func (dbdget *dbDelegate) Rollback() {
	dbdget.dbGorm.Rollback()
}

func (dbdget *dbDelegate) Commit() {
	dbdget.dbGorm.Commit()
}

func connectionstring() string {
	connectionstring := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Name,
		setting.DatabaseSetting.Port)
	fmt.Printf("%s", connectionstring)

	return connectionstring
}