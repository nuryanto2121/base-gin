package contauditlogs

import (
	iauditlogs "app/interface/audit_logs"
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

type contAuditLogs struct {
	useAuditLogs iauditlogs.Usecase
}

func NewContAuditLogs(e *gin.Engine, a iauditlogs.Usecase) {
	controller := &contAuditLogs{
		useAuditLogs: a,
	}

	r := e.Group("/v1/cms/audit_logs")
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
// @Tags AuditLogs
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/audit_logs/{id} [get]
func (c *contAuditLogs) GetDataBy(e *gin.Context) {
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
	data, err := c.useAuditLogs.GetDataBy(ctx, claims, ID)
	if err != nil {
		appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
		return
	}

	appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList AuditLogs
// @Security ApiKeyAuth
// @Tags AuditLogs
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
// @Router /v1/cms/audit_logs [get]
func (c *contAuditLogs) GetList(e *gin.Context) {
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

	responseList, err = c.useAuditLogs.GetList(ctx, claims, paramquery)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "", responseList)
}

// CreateAuditLogs :
// @Summary Add AuditLogs
// @Security ApiKeyAuth
// @Tags AuditLogs
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.AddAuditLogs true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/audit_logs [post]
func (c *contAuditLogs) Create(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		form   models.AddAuditLogs
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

	err = c.useAuditLogs.Create(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusCreated, "Ok", nil)
}

// UpdateAuditLogs :
// @Summary Rubah AuditLogs
// @Security ApiKeyAuth
// @Tags AuditLogs
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Param req body models.AddAuditLogs true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/audit_logs/{id} [put]
func (c *contAuditLogs) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		err    error

		id   = e.Param("id") //kalo bukan int => 0
		form = models.AddAuditLogs{}
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

	// form.UpdatedBy = claims.AuditLogsName
	err = c.useAuditLogs.Update(ctx, claims, SchoolID, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	appE.Response(http.StatusCreated, "Ok", nil)
}

// DeleteAuditLogs :
// @Summary Delete AuditLogs
// @Security ApiKeyAuth
// @Tags AuditLogs
// @Produce  json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Success 200 {object} app.Response
// @Router /v1/cms/audit_logs/{id} [delete]
func (c *contAuditLogs) Delete(e *gin.Context) {
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
	err = c.useAuditLogs.Delete(ctx, claims, ID)
	if err != nil {
		appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}
