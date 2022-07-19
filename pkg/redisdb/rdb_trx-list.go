package redisdb

import (
	"fmt"
	"time"
)

// GetList :
func GetList(key string) ([]string, error) {
	list, err := rdb.SMembers(key).Result()
	return list, err
}

// RemoveList :
func RemoveList(key string, val interface{}) error {
	_, err := rdb.SRem(key, val).Result()
	if err != nil {
		return err
	}
	return nil
}

// AddList :
func AddList(key, val string) error {
	_, err := rdb.SAdd(key, val).Result()
	if err != nil {
		return err
	}
	return nil
}

// TurncateList :
func TurncateList(key string) error {
	_, err := rdb.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

// AddSession :
// func AddSession(key string, val interface{}, mn int) error {
// 	// ss := 1 * time.Hour
// 	var (
// 		tm = time.Minute
// 	)
// 	if mn > 0 {
// 		tm := time.Duration(mn) * time.Minute
// 		fmt.Println(tm)
// 	} else {
// 		tm = 0
// 	}
// 	set := rdb.Set(key, val, tm)
// 	fmt.Println(set)
// 	return nil
// }
func AddSession(key string, val interface{}, mn time.Duration) error {
	set, err := rdb.Set(key, val, mn).Result()
	if err != nil {
		return err
	}
	fmt.Println(set)
	return nil
}

// GetSession :
func GetSession(key string) interface{} {
	value := rdb.Get(key).Val()
	fmt.Println(value)
	return value
}
