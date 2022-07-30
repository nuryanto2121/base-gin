package iusergroup

import (
	"app/models"
	util "app/pkg/utils"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.UserGroup, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result []*models.UserGroup, err error)
	Create(ctx context.Context, data *models.UserGroup) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
	Count(ctx context.Context, queryparam models.ParamList) (result int64, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.UserGroup, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.AddUserGroup) (err error)
	Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddUserGroup) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error)
	DeleteByUserId(ctx context.Context, Claims util.Claims, UserID uuid.UUID) (err error)
}
