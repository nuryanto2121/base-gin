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
	uuid "github.com/satori/go.uuid"
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

// // Update :
// // @Summary Update TermAndConditional
// // @Security ApiKeyAuth
// // @Tags TermAndConditional
// // @Produce json
// // @Param Device-Type header string true "Device Type"
// // @Param Version header string true "Version Apps"
// // @Param Language header string true "Language Apps"
// // @Param req body models.TermAndConditionalForm true "this model set from firebase"
// // @Param id path string true "ID"
// // @Success 200 {object} app.Response
// // @Router /v1/cms/term-and-conditional/{id} [put]
// func (c *contTermAndConditional) Update(e *gin.Context) {
// 	ctx := e.Request.Context()
// 	if ctx == nil {
// 		ctx = context.Background()
// 	}

// 	var (
// 		logger = logging.Logger{}
// 		appE   = app.Gin{C: e}
// 		form   = models.TermAndConditionalForm{}
// 		id     = e.Param("id")
// 	)

// 	Id := uuid.FromStringOrNil(id)
// 	// validasi and bind to struct
// 	httpCode, errMsg := app.BindAndValidMulti(e, &form)
// 	logger.Info(util.Stringify(form))
// 	if httpCode != 200 {
// 		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
// 		return
// 	}

// 	claims, err := app.GetClaims(e)
// 	if err != nil {
// 		appE.ResponseError(tool.GetStatusCode(err), err)
// 		return
// 	}

// 	err = c.useTermAndConditional.Update(ctx, claims, Id, form)
// 	if err != nil {
// 		appE.ResponseError(tool.GetStatusCode(err), err)
// 		return
// 	}
// 	appE.Response(http.StatusOK, "Ok", nil)
// }

// GetById :
// @Summary GetById TermAndConditional
// @Security ApiKeyAuth
// @Tags TermAndConditional
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Success 200 {object} app.Response
// @Router /v1/cms/term-and-conditionalt [get]
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

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	data, err := c.useTermAndConditional.GetDataBy(ctx, claims, Id)
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
// @Router /v1/cms/sku-management [get]
// func (c *contTermAndConditional) GetList(e *gin.Context) {
// 	ctx := e.Request.Context()
// 	if ctx == nil {
// 		ctx = context.Background()
// 	}

// 	var (
// 		logger     = logging.Logger{}
// 		appE       = app.Gin{C: e}
// 		paramquery = models.ParamList{} // ini untuk list
// 		// responseList = models.ResponseModelList{}
// 	)

// 	logger.Info(util.Stringify(paramquery))
// 	httpCode, errMsg := app.BindAndValid(e, &paramquery)
// 	if httpCode != 200 {
// 		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg) // ResponseErrorList(http.StatusBadRequest, errMsg, responseList)
// 		return
// 	}

// 	responseList, err := c.useTermAndConditional.GetList(ctx, paramquery)
// 	if err != nil {
// 		appE.ResponseError(tool.GetStatusCode(err), err)
// 		return
// 	}

// 	appE.Response(http.StatusOK, "Ok", responseList)
// }
