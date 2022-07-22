package contgroup

import (
	igroup "app/interface/group"
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

type ContGroups struct {
	useGroups igroup.Usecase
}

func NewContGroup(e *gin.Engine, useGroups igroup.Usecase) {
	cont := ContGroups{
		useGroups: useGroups,
	}

	r := e.Group("/v1/cms/groups")
	r.POST("", cont.Create)
	r.PUT("/:id", cont.Update)
	r.GET("/:id", cont.GetById)
	r.GET("", cont.GetList)
	r.DELETE("/:id", cont.Delete)
}

// Create :
// @Summary Create Groups
// @Tags Groups
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.GroupForm true "this model set from firebase"
// @Success 200 {object} app.Response
// @Router /v1/cms/groups [post]
func (c *ContGroups) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.GroupForm{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := c.useGroups.Create(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// Update :
// @Summary Update Groups
// @Tags Groups
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.GroupsForm true "this model set from firebase"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/groups/{id} [put]
func (c *ContGroups) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.GroupForm{}
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

	err := c.useGroups.Update(ctx, Id, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// GetById :
// @Summary GetById Groups
// @Tags Groups
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/groups/{id} [get]
func (c *ContGroups) GetById(e *gin.Context) {
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

	data, err := c.useGroups.GetDataBy(ctx, Id)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Groups
// @Tags Groups
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
// @Router /v1/cms/groups [get]
func (c *ContGroups) GetList(e *gin.Context) {
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

	responseList, err := c.useGroups.GetList(ctx, paramquery)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", responseList)
}

// Delete :
// @Summary Delete Groups
// @Tags Groups
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/groups/{id} [delete]
func (c *ContGroups) Delete(e *gin.Context) {
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

	err := c.useGroups.Delete(ctx, Id)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}
