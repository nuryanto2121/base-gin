package ireport

import (
	"app/models"
	"app/pkg/util"
	"context"
)

type Usecase interface {
	GetReport(ctx context.Context, Claims util.Claims, param *models.ReportForm) (interface{}, error)
}

type Repository interface {
	GetReport(ctx context.Context, sWhere, startDate, endDate, userId string) ([]*models.ReportResponse, error)
}
