package itermandconditional

import (
	"app/models"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.TermAndConditional, err error)
	// GetList(ctx context.Context, queryparam models.ParamList) (result []*models.TermAndConditional, err error)
	Create(ctx context.Context, data *models.TermAndConditional) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	// Delete(ctx context.Context, ID uuid.UUID) (err error)
	Count(ctx context.Context, queryparam models.ParamList) (result int64, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.TermAndConditional, err error)
	//GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, data *models.TermAndConditional) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
}
