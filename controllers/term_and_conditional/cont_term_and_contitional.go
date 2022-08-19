package conttermandconditional

import (
	itermandconditional "app/interface/term_and_conditional"
	"app/models"
	"app/pkg/app"
	"app/pkg/logging"
	"app/pkg/middleware"
	tool "app/pkg/tools"
	util "app/pkg/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type contTermAndConditional struct {
	useTermAndConditional itermandconditional.Usecase
}

func NewContTermAndConditional(e *gin.Engine, useTermAndConditional itermandconditional.Usecase) {
	cont := contTermAndConditional{
		useTermAndConditional: useTermAndConditional,
	}

	r := e.Group("/v1/cms/term-and-conditional")
	r.POST("", cont.Create)
	r.Use(middleware.Authorize())
	// r.PUT("/:id", cont.Update)
	r.GET("", cont.GetById)
	// r.GET("", cont.GetList)
	// r.DELETE("/:id", cont.Delete)
}

// Create :
// @Summary Create TermAndConditional
// @Security ApiKeyAuth
// @Tags TermAndConditional
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.TermAndConditionalForm true "this model set from firebase"
// @Success 200 {object} app.Response
// @Router /v1/cms/term-and-conditional [post]
func (c *contTermAndConditional) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.TermAndConditionalForm{}
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

	err = c.useTermAndConditional.Create(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// GetById :
// @Summary GetById TermAndConditional
// @Security ApiKeyAuth
// @Tags TermAndConditional
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Success 200 {object} app.Response
// @Router /v1/cms/term-and-conditional [get]
func (c *contTermAndConditional) GetById(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = app.Gin{C: e}
	)

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	data, err := c.useTermAndConditional.GetDataOne(ctx, claims)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	appE.Response(http.StatusOK, "Ok", data)
}
