package repouserSessionession

import (
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"context"

	iuserSession "app/interface/user_session"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoUserSession struct {
	db db.DBGormDelegate
}

func NewRepoUserSession(Conn db.DBGormDelegate) iuserSession.Repository {
	return &repoUserSession{Conn}
}
func (r *repoUserSession) GetByUser(ctx context.Context, Account uuid.UUID) (result *models.UserSession, err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Where("(email like ? OR phone_no = ?)", Account, Account).First(&result)
	err = query.Error
	if err != nil {
		logger.Error("repo user session GetByUser ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return result, err
}

func (r *repoUserSession) GetByToken(ctx context.Context, Token string) (result *models.UserSession, err error) {
	var (
		sysUser = &models.UserSession{}
		logger  = logging.Logger{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where("token = ? ", Token).Find(sysUser)
	err = query.Error
	if err != nil {
		logger.Error("repo user session GetByToken ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return sysUser, nil
}

func (r *repoUserSession) Create(ctx context.Context, data *models.UserSession) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Create(data)
	err = query.Error
	if err != nil {
		logger.Error("repo user session Create ", err)
		return err
	}
	return nil
}

func (r *repoUserSession) Delete(ctx context.Context, Token string) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Where("token = ?", Token).Delete(&models.UserSession{})
	err = query.Error
	if err != nil {
		logger.Error("repo user session Delete ", err)
		return err
	}
	return nil
}
