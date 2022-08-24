package authorize

import (
	"app/models"
	"app/pkg/db"
	"app/pkg/setting"
	"context"
)

type Session struct {
	Token string `json:"token"`
}

func (s *Session) GetSession(ctx context.Context) (result *models.UserSession, err error) {
	dbConn := db.NewDBdelegate(setting.DatabaseSetting.Debug)
	dbConn.Init()

	conn := dbConn.Get(ctx)
	err = conn.Where("token = ? ", s.Token).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil

}
