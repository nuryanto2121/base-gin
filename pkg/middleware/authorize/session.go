package authorize

import (
	"app/models"
	"context"

	"gorm.io/gorm"
)

type Session struct {
	Token string `json:"token"`
}

func (s *Session) GetSession(ctx context.Context, Conn *gorm.DB) (result *models.UserSession, err error) {
	err = Conn.Where("token = ? ", s.Token).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil

}
