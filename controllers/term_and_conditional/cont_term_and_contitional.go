package conttermandconditional

import (
	itermandconditional "app/interface/term_and_conditional"
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

type contTermAndConditional struct {
	useTermAndConditional itermandconditional.Usecase
}

func NewContTermAndConditional(e *gin.Engine, useTermAndConditional itermandconditional.Usecase) {
	cont := contTermAndConditional{
		useTermAndConditional: useTermAndConditional,
	}

	r := e.Group("/v1/cms/termandconditional")
	r.POST("", cont.Create)
	r.PUT("/:id", cont.Update)
	r.GET("/:id", cont.GetById)
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
// @Param req body models.TermAndConditional true "this model set from firebase"
// @Success 200 {object} app.Response
// @Router /v1/cms/termandconditional [post]
func (c *contTermAndConditional) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.TermAndConditional{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := c.useTermAndConditional.Create(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// Update :
// @Summary Update TermAndConditional
// @Security ApiKeyAuth
// @Tags TermAndConditional
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.TermAndConditional true "this model set from firebase"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/termandconditional/{id} [put]
func (c *contTermAndConditional) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.TermAndConditional{}
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

	err := c.useTermAndConditional.Update(ctx, Id, &form)
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
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/termandconditionalt/{id} [get]
func (c *contTermAndConditional) GetById(e *gin.Context) {
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

	data, err := c.useTermAndConditional.GetDataBy(ctx, Id)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}
