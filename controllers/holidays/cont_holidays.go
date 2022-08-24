package contholidays

import (
	iholidays "app/interface/holidays"
	"app/models"
	"app/pkg/app"
	"app/pkg/logging"
	"app/pkg/middleware"
	tool "app/pkg/tools"
	util "app/pkg/util"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type ContHolidays struct {
	useHolidays iholidays.Usecase
}

func NewContHolidays(e *gin.Engine, useHolidays iholidays.Usecase) {
	cont := ContHolidays{
		useHolidays: useHolidays,
	}

	r := e.Group("/v1/cms/holidays")
	r.Use(middleware.Authorize())
	r.POST("", cont.Create)
	r.PUT("/:id", cont.Update)
	r.GET("/:id", cont.GetById)
	r.GET("", cont.GetList)
	r.DELETE("/:id", cont.Delete)
}

// Create :
// @Summary Create Holidays
// @Tags Holidays
// @Security ApiKeyAuth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.HolidayForm true "this model set from firebase"
// @Success 200 {object} app.Response
// @Router /v1/cms/holidays [post]
func (c *ContHolidays) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.HolidayForm{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	err = c.useHolidays.Create(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "Ok", nil)
}

// Update :
// @Summary Update Holidays
// @Tags Holidays
// @Security ApiKeyAuth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Param req body models.HolidayForm true "this model set from firebase"
// @Success 200 {object} app.Response
// @Router /v1/cms/holidays/{id} [put]
func (c *ContHolidays) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = &models.HolidayForm{}
		id     = e.Param("id")
	)

	Id := uuid.FromStringOrNil(id)
	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	err = c.useHolidays.Update(ctx, claims, Id, form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// GetById :
// @Summary GetById Holidays
// @Tags Holidays
// @Security ApiKeyAuth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/holidays/{id} [get]
func (c *ContHolidays) GetById(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		id     = e.Param("id")
	)

	Id := uuid.FromStringOrNil(id)
	logger.Info(id)

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	data, err := c.useHolidays.GetDataBy(ctx, claims, Id)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Holidays
// @Tags Holidays
// @Security ApiKeyAuth
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param page query int true "Page"
// @Param perpage query int true "PerPage"
// @Param search query string false "Search"
// @Param initsearch query string false "InitSearch"
// @Param sortfield query string false "SortField"
// @Success 200 {object} models.ResponseModelList
// @Router /v1/cms/holidays [get]
func (c *ContHolidays) GetList(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger     = logging.Logger{}
		appE       = app.Gin{C: e}
		paramquery = models.ParamList{} // ini untuk list
		// responseList = models.ResponseModelList{}
	)

	logger.Info(util.Stringify(paramquery))
	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg) // ResponseErrorList(http.StatusBadRequest, errMsg, responseList)
		return
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	responseList, err := c.useHolidays.GetList(ctx, claims, paramquery)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", responseList)
}

// Delete :
// @Summary Delete Holidays
// @Tags Holidays
// @Security ApiKeyAuth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/holidays/{id} [delete]
func (c *ContHolidays) Delete(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		id     = e.Param("id")
	)

	Id := uuid.FromStringOrNil(id)
	logger.Info(id)

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	err = c.useHolidays.Delete(ctx, claims, Id)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}
