package continventory

import (
	iinventory "app/interface/inventory"
	iorder "app/interface/order"
	"app/models"
	"app/pkg/middleware"
	"context"
	"net/http"

	//app "app/pkg"
	"app/pkg/app"
	"app/pkg/logging"
	tool "app/pkg/tools"
	util "app/pkg/util"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type contInventory struct {
	useInventory iinventory.Usecase
	useOrder     iorder.Usecase
}

func NewContInventory(e *gin.Engine, a iinventory.Usecase, b iorder.Usecase) {
	controller := &contInventory{
		useInventory: a,
		useOrder:     b,
	}

	r := e.Group("/v1/cms/inventory")
	r.Use(middleware.Authorize())
	r.Use(middleware.Versioning())
	r.POST("/:id", controller.Update)
	r.POST("/status", controller.Status)

}

// UpdateOrSaveInventory :
// @Summary Rubah atau simpan Inventory
// @Security ApiKeyAuth
// @Tags Inventory
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param id path string true "ID"
// @Param req body models.InventoryForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/inventory/{id} [post]
func (c *contInventory) Update(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		err    error

		id   = e.Param("id") //kalo bukan int => 0
		form = models.InventoryForm{}
	)

	ID, err := uuid.FromString(id)

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

	// form.UpdatedBy = claims.InventoryName
	err = c.useInventory.Save(ctx, claims, ID, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	appE.Response(http.StatusCreated, "Ok", nil)
}

// Status :
// @Summary Approve or Reject order inventory
// @Security ApiKeyAuth
// @Tags Inventory
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.InventoryStatusForm true "req param #status int64 0 = SUBMITTED , 1 = APPROVE, 2 = Reject"
// @Success 200 {object} app.Response
// @Router /v1/cms/inventory/status [post]
func (c *contInventory) Status(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		err    error

		form = models.InventoryStatusForm{}
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

	err = c.useOrder.UpdateStatus(ctx, claims, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}
	appE.Response(http.StatusCreated, "Ok", nil)
}
