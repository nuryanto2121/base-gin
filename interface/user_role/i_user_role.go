package iuserrole

import (
	"app/models"
	util "app/pkg/utils"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetDataBy(ctx context.Context, key, value string) (result *models.UserRoleDesc, err error)
	GetById(ctx context.Context, ID uuid.UUID) (result *models.UserRole, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result []*models.UserRole, err error)
	GetListByUser(ctx context.Context, key, value string) (result []*models.UserRoleDesc, err error)
	Create(ctx context.Context, data *models.UserRole) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
	Count(ctx context.Context, queryparam models.ParamList) (result int64, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, key, value string) (result *models.UserRoleDesc, err error)
	GetById(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.UserRole, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.AddUserRole) (err error)
	Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddUserRole) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error)
	DeleteByUserId(ctx context.Context, Claims util.Claims, UserID uuid.UUID) (err error)
}
