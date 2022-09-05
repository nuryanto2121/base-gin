package conttransaction

import (
	itransaction "app/interface/transaction"
	"app/models"
	"app/pkg/middleware"
	"context"
	"fmt"
	"net/http"

	//app "app/pkg"
	"app/pkg/app"
	"app/pkg/logging"
	tool "app/pkg/tools"
	"app/pkg/util"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type contTransaction struct {
	useTransaction itransaction.Usecase
}

func NewContTransaction(e *gin.Engine, a itransaction.Usecase) {
	controller := &contTransaction{
		useTransaction: a,
	}

	r := e.Group("/v1/transaction")
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
// @Tags Transaction
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/transaction/{id} [get]
func (c *contTransaction) GetDataBy(e *gin.Context) {
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
		appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
		return
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
		return
	}
	data, err := c.useTransaction.GetDataBy(ctx, claims, ID)
	if err != nil {
		appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Transaction
// @Security ApiKeyAuth
// @Tags Transaction
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param page query int true "Page"
// @Param perpage query int true "PerPage"
// @Param search query string false "Search"
// @Param initsearch query string false "InitSearch | status_transaction => STATUS_ORDER=2000001 |STATUS_CHECKIN=2000002|STATUS_CHECKOUT=2000003"
// @Param sortfield query string false "SortField"
// @Success 200 {object} models.ResponseModelList
// @Router /v1/transaction [get]
func (c *contTransaction) GetList(e *gin.Context) {
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

	responseList, err = c.useTransaction.GetList(ctx, claims, paramquery)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "", responseList)
}

// CreateTransaction :
// @Summary Add Transaction
// @Security ApiKeyAuth
// @Tags Transaction
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.TransactionForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/transaction [post]
func (c *contTransaction) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		form   models.TransactionForm
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

	err = c.useTransaction.Create(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "Ok", nil)
}

// UpdateTransaction :
// @Summary Rubah Transaction
// @Security ApiKeyAuth
// @Tags Transaction
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Param req body models.TransactionForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/transaction/{id} [put]
func (c *contTransaction) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		err    error

		id   = e.Param("id") //kalo bukan int => 0
		form = models.TransactionForm{}
	)

	ID, err := uuid.FromString(id)
	logger.Info(ID)
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

	// form.UpdatedBy = claims.TransactionName
	err = c.useTransaction.Update(ctx, claims, ID, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	appE.Response(http.StatusCreated, "Ok", nil)
}

// DeleteTransaction :
// @Summary Delete Transaction
// @Security ApiKeyAuth
// @Tags Transaction
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/transaction/{id} [delete]
func (c *contTransaction) Delete(e *gin.Context) {
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
		appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
		return
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	err = c.useTransaction.Delete(ctx, claims, ID)
	if err != nil {
		appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}
