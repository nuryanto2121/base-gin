package iinventory

import (
	"app/models"
	util "app/pkg/utils"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Inventory, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Inventory, err error)
	Create(ctx context.Context, data *models.Inventory) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
	Count(ctx context.Context, queryparam models.ParamList) (result int64, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.Inventory, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.InventoryForm) (err error)
	Save(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.InventoryForm) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error)
}
