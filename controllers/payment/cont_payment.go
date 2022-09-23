package payment

import (
	ipayment "app/interface/payment"
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

type contPayment struct {
	usePayment ipayment.Usecase
}

func NewContPayment(e *gin.Engine, a ipayment.Usecase) {
	controller := &contPayment{
		usePayment: a,
	}
	r := e.Group("/v1/payment")

	r.POST("/receive", controller.Receive)
	r.Use(middleware.Authorize()).
		Use(middleware.Versioning()).
		POST("", controller.Payment)
}

// CreatePayment :
// @Summary Receive Payment
// @Security ApiKeyAuth
// @Tags Payment
// @Produce json
// @Param req body models.MidtransNotification true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/payment/receive [post]
func (c *contPayment) Receive(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		form   models.MidtransNotification
	)

	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := c.usePayment.Receive(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "", nil)
}

// CreatePayment :
// @Summary Add Payment from mobile
// @Security ApiKeyAuth
// @Tags Payment
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.TransactionPaymentForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/payment [post]
func (c *contPayment) Payment(e *gin.Context) {
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

	rest, err := c.usePayment.Payment(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "payment success", rest)
}
