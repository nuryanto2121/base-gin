package redisdb

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/369-engineer/369backend/account/pkg/setting"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

// Setup :
func Setup() {
	now := time.Now()
	conString := fmt.Sprintf("%s:%d", setting.RedisDBSetting.Host, setting.RedisDBSetting.Port)
	rdb = redis.NewClient(&redis.Options{
		Addr:     conString,
		Password: setting.RedisDBSetting.Password,
		DB:       setting.RedisDBSetting.DB,
	})

	fmt.Printf("\nconnection String redis : %s\n", conString)
	_, err := rdb.Ping().Result()
	if err != nil {
		fmt.Println(err)
		panic(err)

	}

	timeSpent := time.Since(now)
	log.Printf("Config redis is ready in %v", timeSpent)
}
