package iskumanagement

import (
	"app/models"
	util "app/pkg/utils"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetDataBy(ctx context.Context, key, value string) (result *models.SkuManagement, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result []*models.SkuManagement, err error)
	Create(ctx context.Context, data *models.SkuManagement) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
	Count(ctx context.Context, queryparam models.ParamList) (result int64, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.SkuManagement, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.AddSkuManagement) (err error)
	Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddSkuManagement) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error)
}
