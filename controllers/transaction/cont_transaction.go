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

	r := e.Group("/v1/cms/transaction")
	r.Use(middleware.Authorize())
	r.Use(middleware.Versioning())
	r.GET("/scan", controller.GetDataBy)
	r.GET("", controller.GetList)
	r.POST("/payment", controller.Payment)
	r.GET("/print-ticket", controller.GetTicket)
	r.POST("/check-in", controller.CheckIn)
	r.POST("/check-out", controller.CheckOut)

	l := e.Group("v1/transaction")
	l.Use(middleware.Authorize())
	l.Use(middleware.Versioning())
	l.POST("", controller.Create)
	l.DELETE("/:id", controller.Delete)
	l.PUT("/:id", controller.Update)
	l.GET("/tickets", controller.GetListTicket)
}

// GetDataByID :
// @Summary GetById for scan ticket / detail transaction
// @Security ApiKeyAuth
// @Tags Transaction
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param transactionCode query string true "transactionCode"
// @Success 200 {object} app.Response
// @Router /v1/cms/transaction/scan [get]
func (c *contTransaction) GetDataBy(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger          = logging.Logger{}
		appE            = app.Gin{C: e} // wajib
		transactionCode = e.Query("transactionCode")
	)
	logger.Info(transactionCode)

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
		return
	}
	data, err := c.useTransaction.GetDataBy(ctx, claims, transactionCode)
	if err != nil {
		appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}

// GetTicket :
// @Summary GetById for scan ticket / detail transaction
// @Security ApiKeyAuth
// @Tags Transaction
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param transactionId query string true "transactionId"
// @Success 200 {object} app.Response
// @Router /v1/cms/transaction/print-ticket [get]
func (c *contTransaction) GetTicket(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger        = logging.Logger{}
		appE          = app.Gin{C: e}            // wajib
		transactionId = e.Query("transactionId") //e.Param("transactionId") //kalo bukan int => 0
	)
	logger.Info(transactionId)

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
		return
	}
	claims.Role = "ticket"
	data, err := c.useTransaction.GetDataBy(ctx, claims, transactionId)
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
// @Router /v1/cms/transaction [get]
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

// GetListTicket :
// @Summary GetListTicket Transaction
// @Security ApiKeyAuth
// @Tags Transaction
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
// @Router /v1/transaction/tickets [get]
func (c *contTransaction) GetListTicket(e *gin.Context) {
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

	responseList, err = c.useTransaction.GetListTicketUser(ctx, claims, paramquery)
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

	response, err := c.useTransaction.Create(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "Ok", response)
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

// CreateTransaction :
// @Summary Add Payment CMS
// @Security ApiKeyAuth
// @Tags Transaction
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.TransactionPaymentForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/transaction/payment [post]
func (c *contTransaction) Payment(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		form   models.TransactionPaymentForm
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

	err = c.useTransaction.Payment(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "payment success", nil)
}

// CheckInTransaction :
// @Summary Check in transaction
// @Security ApiKeyAuth
// @Tags Transaction
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.CheckInCheckOutForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/transaction/check-in [post]
func (c *contTransaction) CheckIn(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		form   models.CheckInCheckOutForm
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

	err = c.useTransaction.CheckIn(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "Ok", nil)
}

// CheckOutTransaction :
// @Summary Check out transaction
// @Security ApiKeyAuth
// @Tags Transaction
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.CheckInCheckOutForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/transaction/check-out [post]
func (c *contTransaction) CheckOut(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		form   models.CheckInCheckOutForm
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

	err = c.useTransaction.CheckOut(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "Ok", nil)
}
