package iuserapps

import (
	"app/models"
	"app/pkg/util"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetDataBy(ctx context.Context, key, value string) (result *models.UserApps, err error)
	GetByAccount(ctx context.Context, account string) (result *models.UserApps, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result []*models.UserApps, err error)
	Create(ctx context.Context, data *models.UserApps) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
	Count(ctx context.Context, queryparam models.ParamList) (result int64, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.UserApps, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.UserApps) (err error)
	UpsertChild(ctx context.Context, Claims util.Claims, data models.ChildForm) (result models.ChildForm, err error)
	Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddUserApps) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error)
}
