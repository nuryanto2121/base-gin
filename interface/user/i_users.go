package iusers

import (
	"context"

	"app/models"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetById(ctx context.Context, ID uuid.UUID) (result *models.Users, err error)
	GetDataBy(ctx context.Context, key, value string) (result *models.Users, err error)
	IsExist(ctx context.Context, key, value string) (bool, error)
	GetByAccount(ctx context.Context, Account string) (result *models.Users, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result []*models.ListUserCms, err error)
	Create(ctx context.Context, data *models.Users) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
	Count(ctx context.Context, queryparam models.ParamList) (result int64, err error)
}
type Usecase interface {
	GetDataBy(ctx context.Context, ID uuid.UUID) (result interface{}, err error)
	GetByEmailSaUser(ctx context.Context, email string) (result *models.Users, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error)
	CreateCms(ctx context.Context, data *models.AddUserCms) (err error)
	Update(ctx context.Context, ID uuid.UUID, data *models.EditUserCms) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
}
