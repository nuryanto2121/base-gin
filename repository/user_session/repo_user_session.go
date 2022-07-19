package repouserSessionession

import (
	"app/models"
	"context"

	iuserSession "app/interface/user_session"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoUserSession struct {
	Conn *gorm.DB
}

func NewRepoUserSession(Conn *gorm.DB) iuserSession.Repository {
	return &repoUserSession{Conn}
}
func (db *repoUserSession) GetByUser(ctx context.Context, Account uuid.UUID) (result *models.UserSession, err error) {

	query := db.Conn.WithContext(ctx).Where("(email like ? OR phone_no = ?)", Account, Account).First(&result)
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return result, err
}

func (db *repoUserSession) GetByToken(ctx context.Context, Token string) (result *models.UserSession, err error) {
	var sysUser = &models.UserSession{}
	query := db.Conn.WithContext(ctx).Where("token = ? ", Token).Find(sysUser)
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return sysUser, nil
}

func (db *repoUserSession) Create(ctx context.Context, data *models.UserSession) (err error) {
	query := db.Conn.WithContext(ctx).Create(data)
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoUserSession) Delete(ctx context.Context, Token string) (err error) {

	query := db.Conn.WithContext(ctx).Where("token = ?", Token).Delete(&models.UserSession{})
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
