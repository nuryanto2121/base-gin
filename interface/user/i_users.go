package iusers

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/369-engineer/369backend/account/models"
)

type Repository interface {
	GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Users, err error)
	GetByAccount(ctx context.Context, Account string) (result *models.Users, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Users, err error)
	Create(ctx context.Context, data *models.Users) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
	Count(ctx context.Context, queryparam models.ParamList) (result int64, err error)
}
type Usecase interface {
	GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Users, err error)
	GetByEmailSaUser(ctx context.Context, email string) (result *models.Users, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, data *models.Users) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
}
