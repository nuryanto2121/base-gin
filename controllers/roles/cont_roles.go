package contgroup

import (
	irole "app/interface/role"
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

type ContRoles struct {
	useRoles irole.Usecase
}

func NewContRole(e *gin.Engine, useRoles irole.Usecase) {
	cont := ContRoles{
		useRoles: useRoles,
	}

	r := e.Group("/v1/cms/roles")
	r.Use(middleware.Authorize())
	r.POST("", cont.Create)
	r.PUT("/:id", cont.Update)
	r.GET("/:id", cont.GetById)
	r.GET("", cont.GetList)
	r.DELETE("/:id", cont.Delete)
}

// Create :
// @Summary Create Roles
// @Tags Roles
// @Security ApiKeyAuth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.RoleForm true "this model set from firebase"
// @Success 200 {object} app.Response
// @Router /v1/cms/roles [post]
func (c *ContRoles) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.RoleForm{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := c.useRoles.Create(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// Update :
// @Summary Update Roles
// @Tags Roles
// @Security ApiKeyAuth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.RoleForm true "this model set from firebase"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/roles/{id} [put]
func (c *ContRoles) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.RoleForm{}
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

	err := c.useRoles.Update(ctx, Id, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// GetById :
// @Summary GetById Roles
// @Tags Roles
// @Security ApiKeyAuth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/roles/{id} [get]
func (c *ContRoles) GetById(e *gin.Context) {
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

	data, err := c.useRoles.GetDataBy(ctx, Id)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Roles
// @Tags Roles
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
// @Router /v1/cms/roles [get]
func (c *ContRoles) GetList(e *gin.Context) {
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

	responseList, err := c.useRoles.GetList(ctx, paramquery)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", responseList)
}

// Delete :
// @Summary Delete Roles
// @Tags Roles
// @Security ApiKeyAuth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/roles/{id} [delete]
func (c *ContRoles) Delete(e *gin.Context) {
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

	err := c.useRoles.Delete(ctx, Id)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}
