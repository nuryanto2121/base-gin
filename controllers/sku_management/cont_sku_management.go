package contskumanagement

import (
	iskumanagement "app/interface/sku_management"
	"app/models"
	"app/pkg/app"
	"app/pkg/logging"
	tool "app/pkg/tools"
	util "app/pkg/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type contskumanagement struct {
	useskumanagement iskumanagement.Usecase
}

func NewContSkuManagement(e *gin.Engine, useskumanagement iskumanagement.Usecase) {
	cont := contskumanagement{
		useskumanagement: useskumanagement,
	}

	r := e.Group("/v1/cms/skumanagement")
	r.POST("", cont.Create)
	r.PUT("/:id", cont.Update)
	r.GET("/:id", cont.GetById)
	r.GET("", cont.GetList)
	r.DELETE("/:id", cont.Delete)
}

// Create :
// @Summary Create SkuManagement
// @Security ApiKeyAuth
// @Tags SkuManagement
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.SkuManagement true "this model set from firebase"
// @Success 200 {object} app.Response
// @Router /v1/cms/skumanagement [post]
func (c *contskumanagement) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.SkuManagement{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := c.useskumanagement.Create(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// Update :
// @Summary Update SkuManagement
// @Security ApiKeyAuth
// @Tags SkuManagement
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.SkuManagement true "this model set from firebase"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/skumanagement/{id} [put]
func (c *contskumanagement) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.SkuManagement{}
		id     = e.Param("id")
	)

	Id := uuid.FromStringOrNil(id)
	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := c.useskumanagement.Update(ctx, Id, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// GetById :
// @Summary GetById SkuManagement
// @Security ApiKeyAuth
// @Tags SkuManagement
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/skumanagement/{id} [get]
func (c *contskumanagement) GetById(e *gin.Context) {
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

	data, err := c.useskumanagement.GetDataBy(ctx, Id)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList SkuManagement
// @Security ApiKeyAuth
// @Tags SkuManagement
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
// @Router /v1/cms/skumanagement [get]
func (c *contskumanagement) GetList(e *gin.Context) {
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

	responseList, err := c.useskumanagement.GetList(ctx, paramquery)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", responseList)
}

// Delete :
// @Summary Delete SKuManagement
// @Security ApiKeyAuth
// @Tags SkuManagement
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/skumanagement/{id} [delete]
func (c *contskumanagement) Delete(e *gin.Context) {
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

	err := c.useskumanagement.Delete(ctx, Id)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}