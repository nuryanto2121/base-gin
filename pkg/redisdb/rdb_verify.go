package redisdb

import (
	"encoding/json"

	"gitlab.com/369-engineer/369backend/account/pkg/setting"

	"github.com/mitchellh/mapstructure"
)

// Verify :
type Verify struct {
	Email      string
	UserName   string
	VerifyLink string
}

// StoreVerify :
func StoreVerify(data interface{}) error {
	var verify Verify

	err := mapstructure.Decode(data, &verify)
	if err != nil {
		return err
	}

	mVerify := map[string]interface{}{
		"mail_type": "verify",
		"data":      verify,
	}

	dVerify, err := json.Marshal(mVerify)
	if err != nil {
		return err
	}

	_, err = rdb.SAdd(setting.RedisDBSetting.Key, string(dVerify)).Result()
	if err != nil {
		return err
	}

	return nil
}
