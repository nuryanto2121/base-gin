package continventory

import (
	iinventory "app/interface/inventory"
	"app/models"
	"app/pkg/middleware"
	"context"
	"net/http"

	//app "app/pkg"
	"app/pkg/app"
	"app/pkg/logging"
	tool "app/pkg/tools"
	util "app/pkg/utils"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type contInventory struct {
	useInventory iinventory.Usecase
}

func NewContInventory(e *gin.Engine, a iinventory.Usecase) {
	controller := &contInventory{
		useInventory: a,
	}

	r := e.Group("/v1/cms/inventory")
	r.Use(middleware.Authorize())
	//r.Use(midd.Versioning)
	r.POST("/:id", controller.Update)

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
// @Param req body models.AddInventory true "req param #changes are possible to adjust the form of the registration form from frontend"
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
		form = models.AddInventory{}
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