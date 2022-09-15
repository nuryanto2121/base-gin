package contorder

import (
	iorder "app/interface/order"
	"app/models"
	"app/pkg/middleware"
	"context"
	"net/http"

	//app "app/pkg"
	"app/pkg/app"
	"app/pkg/logging"
	tool "app/pkg/tools"
	util "app/pkg/util"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type contOrder struct {
	useOrder iorder.Usecase
}

func NewContOrder(e *gin.Engine, a iorder.Usecase) {
	controller := &contOrder{
		useOrder: a,
	}

	r := e.Group("/v1/cms/order")
	r.Use(middleware.Authorize())
	r.Use(middleware.Versioning())
	r.GET("/:id", controller.GetDataBy)
	r.GET("", controller.GetList)
	r.POST("", controller.Create)
	r.PUT("/:id", controller.Update)
	r.DELETE("/:id", controller.Delete)
}

// GetDataByID :
// @Summary GetById
// @Security ApiKeyAuth
// @Tags Order
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/order/{id} [get]
func (c *contOrder) GetDataBy(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e} // wajib
		id     = e.Param("id") //kalo bukan int => 0
	)
	ID, err := uuid.FromString(id)
	logger.Info(ID)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	data, err := c.useOrder.GetDataBy(ctx, claims, ID)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Order
// @Security ApiKeyAuth
// @Tags Order
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
// @Router /v1/cms/order [get]
func (c *contOrder) GetList(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger       = logging.Logger{}
		appE         = app.Gin{C: e}      // wajib
		paramquery   = models.ParamList{} // ini untuk list
		responseList = models.ResponseModelList{}
		err          error
	)

	httpCode, errMsg := app.BindAndValidMulti(e, &paramquery)
	logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	responseList, err = c.useOrder.GetList(ctx, claims, paramquery)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "", responseList)
}

// CreateOrder :
// @Summary Add Order
// @Security ApiKeyAuth
// @Tags Order
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.AddOrder true "req param #status int64 0 = SUBMITTED , 1 = APPROVE"
// @Success 200 {object} app.Response
// @Router /v1/cms/order [post]
func (c *contOrder) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		form   models.AddOrder
	)

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

	err = c.useOrder.Create(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "Ok", nil)
}

// UpdateOrder :
// @Summary Rubah Order
// @Security ApiKeyAuth
// @Tags Order
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Param req body models.AddOrder true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/order/{id} [put]
func (c *contOrder) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		err    error

		id   = e.Param("id") //kalo bukan int => 0
		form = models.AddOrder{}
	)

	SchoolID, err := uuid.FromString(id)
	logger.Info(SchoolID)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

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

	// form.UpdatedBy = claims.OrderName
	err = c.useOrder.Update(ctx, claims, SchoolID, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	appE.Response(http.StatusCreated, "Ok", nil)
}

// DeleteOrder :
// @Summary Delete Order
// @Security ApiKeyAuth
// @Tags Order
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/order/{id} [delete]
func (c *contOrder) Delete(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e} // wajib
		id     = e.Param("id")
	)
	ID, err := uuid.FromString(id)
	logger.Info(ID)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	err = c.useOrder.Delete(ctx, claims, ID)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}
