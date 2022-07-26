package contoutlets

import (
	ioutlets "app/interface/outlets"
	"app/models"
	"app/pkg/middleware"
	"context"
	"fmt"
	"net/http"

	//app "app/pkg"
	"app/pkg/app"
	"app/pkg/logging"
	tool "app/pkg/tools"
	util "app/pkg/utils"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type contOutlets struct {
	useOutlets ioutlets.Usecase
}

func NewContOutlets(e *gin.Engine, a ioutlets.Usecase) {
	controller := &contOutlets{
		useOutlets: a,
	}

	r := e.Group("/v1/cms/outlets")
	r.Use(middleware.Authorize())
	//r.Use(midd.Versioning)
	r.GET("/:id", controller.GetDataBy)
	r.GET("", controller.GetList)
	r.POST("", controller.Create)
	r.PUT("/:id", controller.Update)
	r.DELETE("/:id", controller.Delete)
}

// GetDataByID :
// @Summary GetById
// @Security ApiKeyAuth
// @Tags Outlets
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/outlets/{id} [get]
func (c *contOutlets) GetDataBy(e *gin.Context) {
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
	data, err := c.useOutlets.GetDataBy(ctx, claims, ID)
	if err != nil {
		appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Outlets
// @Security ApiKeyAuth
// @Tags Outlets
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
// @Router /v1/cms/outlets [get]
func (c *contOutlets) GetList(e *gin.Context) {
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

	responseList, err = c.useOutlets.GetList(ctx, claims, paramquery)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "", responseList)
}

// CreateOutlets :
// @Summary Add Outlets
// @Security ApiKeyAuth
// @Tags Outlets
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.OutletForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/outlets [post]
func (c *contOutlets) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		form   models.OutletForm
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

	err = c.useOutlets.Create(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "Ok", nil)
}

// UpdateOutlets :
// @Summary Rubah Outlets
// @Security ApiKeyAuth
// @Tags Outlets
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Param req body models.AddOutlets true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/outlets/{id} [put]
func (c *contOutlets) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		err    error

		id   = e.Param("id") //kalo bukan int => 0
		form = models.AddOutlets{}
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

	// form.UpdatedBy = claims.OutletsName
	err = c.useOutlets.Update(ctx, claims, SchoolID, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	appE.Response(http.StatusCreated, "Ok", nil)
}

// DeleteOutlets :
// @Summary Delete Outlets
// @Security ApiKeyAuth
// @Tags Outlets
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/outlets/{id} [delete]
func (c *contOutlets) Delete(e *gin.Context) {
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
	err = c.useOutlets.Delete(ctx, claims, ID)
	if err != nil {
		appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}
