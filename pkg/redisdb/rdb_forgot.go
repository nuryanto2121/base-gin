package redisdb

import (
	"encoding/json"

	"app/pkg/setting"

	"github.com/mitchellh/mapstructure"
)

// Forgot :
type Forgot struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	ButtonLink string `json:"button_link"`
}

// StoreForgot :
func StoreForgot(data interface{}) error {
	var forgot Forgot

	err := mapstructure.Decode(data, &forgot)
	if err != nil {
		return err
	}

	bForgot, err := json.Marshal(forgot)
	if err != nil {
		return err
	}

	mForgot := map[string]interface{}{
		"email_type": "forgot",
		"data":       string(bForgot),
	}

	dForgot, _ := json.Marshal(mForgot)
	if err != nil {
		return err
	}

	_, err = rdb.SAdd(setting.RedisDBSetting.Key, string(dForgot)).Result()
	if err != nil {
		return err
	}

	return nil
}
