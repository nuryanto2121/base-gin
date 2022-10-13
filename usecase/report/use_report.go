package usereport

import (
	ireport "app/interface/report"
	iskumanagement "app/interface/sku_management"
	"app/models"
	"app/pkg/util"
	"context"
	"fmt"
	"time"
)

type useReport struct {
	repoReport     ireport.Repository
	repoSku        iskumanagement.Repository
	contextTimeOut time.Duration
}

func NewUseReport(repoReport ireport.Repository, repoSku iskumanagement.Repository, timeout time.Duration) ireport.Usecase {
	return &useReport{
		repoReport:     repoReport,
		repoSku:        repoSku,
		contextTimeOut: timeout,
	}
}

// GetReport implements ireport.Usecase
func (u *useReport) GetReport(ctx context.Context, Claims util.Claims, param *models.ReportForm) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		userId        = Claims.UserID
		outletFIleter = ""
	)
	if param.OutletId != "" {
		outletFIleter = fmt.Sprintf(` AND o.id = '%s'`, param.OutletId)
	}

	startDate := param.StartDate.Format("2006-01-02")
	endDate := param.EndDate.Format("2006-01-02")

	report, err := u.repoReport.GetReport(ctx, outletFIleter, startDate, endDate, userId)
	if err != nil {
		return nil, err
	}

	// sku, err := u.repoSku.GetList(ctx, models.ParamList{
	// 	Page:      1,
	// 	PerPage:   10000,
	// 	SortField: "case when is_bracelet =true then 1 else 2 end,sku_name",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	response := map[string]interface{}{
		"detail": report,
		// "sku":    sku,
	}
	return response, nil
}
