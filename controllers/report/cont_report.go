package contreport

import (
	ireport "app/interface/report"
	"app/models"
	"app/pkg/app"
	"app/pkg/logging"
	"app/pkg/middleware"
	tool "app/pkg/tools"
	"app/pkg/util"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContReport struct {
	useReport ireport.Usecase
}

func NewContReport(e *gin.Engine, useReport ireport.Usecase) {
	cont := ContReport{
		useReport: useReport,
	}
	r := e.Group("/v1/cms/report")
	r.Use(middleware.Authorize())
	r.Use(middleware.Versioning())
	r.GET("", cont.GetReport)
}

// GetList :
// @Summary GetList Report
// @Tags Report
// @Security ApiKeyAuth
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param start_date query string true "StartDate"
// @Param end_date query string true "EndDate"
// @Param outlet_id query string false "OutletId"
// @Success 200 {object} models.ResponseModelList
// @Router /v1/cms/report [get]
func (c *ContReport) GetReport(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger     = logging.Logger{}
		appE       = app.Gin{C: e}
		paramquery = models.ReportForm{} // ini untuk list
		// responseList = models.ResponseModelList{}
	)

	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg) // ResponseErrorList(http.StatusBadRequest, errMsg, responseList)
		return
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	responseList, err := c.useReport.GetReport(ctx, claims, &paramquery)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", responseList)
}
