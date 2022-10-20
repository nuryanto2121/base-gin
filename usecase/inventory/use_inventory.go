package useinventory

import (
	iauditlogs "app/interface/audit_logs"
	iinventory "app/interface/inventory"
	ioutlets "app/interface/outlets"
	iskumanagement "app/interface/sku_management"
	"app/models"
	"app/pkg/logging"
	util "app/pkg/util"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jinzhu/copier"

	uuid "github.com/satori/go.uuid"
)

type useInventory struct {
	repoInventory  iinventory.Repository
	useAuditLogs   iauditlogs.Usecase
	repoOutlet     ioutlets.Repository
	repoProduct    iskumanagement.Repository
	contextTimeOut time.Duration
}

func NewUseInventory(
	repoInventory iinventory.Repository,
	useAuditLogs iauditlogs.Usecase,
	repoOutlet ioutlets.Repository,
	repoProduct iskumanagement.Repository,
	timeout time.Duration,
) iinventory.Usecase {
	return &useInventory{
		repoInventory:  repoInventory,
		useAuditLogs:   useAuditLogs,
		repoOutlet:     repoOutlet,
		repoProduct:    repoProduct,
		contextTimeOut: timeout,
	}
}

func (u *useInventory) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.Inventory, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoInventory.GetDataBy(ctx, ID)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useInventory) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoInventory.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoInventory.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useInventory) Create(ctx context.Context, Claims util.Claims, data *models.InventoryForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mInventory = models.Inventory{}
	)

	// mapping to struct model saRole
	err = copier.Copy(&mInventory.AddInventory, data)
	if err != nil {
		return err
	}

	mInventory.Qty = data.QtyChange
	mInventory.CreatedBy = uuid.FromStringOrNil(Claims.UserID)
	mInventory.UpdatedBy = uuid.FromStringOrNil(Claims.UserID)

	err = u.repoInventory.Create(ctx, &mInventory)
	if err != nil {
		return err
	}
	return nil

}

func (u *useInventory) Save(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.InventoryForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// var (
	// 	qty       int64 = 0
	// 	qtyChange int64 = 0
	// 	qtyDelta  int64 = 0
	// 	auditLogs       = &models.AddAuditLogs{
	// 		AuditDate: util.GetTimeNow(),
	// 		Source:    "outlet-inventory",
	// 		Username:  Claims.UserName,
	// 	}
	// )

	if ID == uuid.Nil {

		//save data
		err := u.Create(ctx, Claims, data)
		if err != nil {
			return err
		}

	} else {
		//getDataOld
		invOld, err := u.repoInventory.GetDataBy(ctx, ID)
		if err != nil {
			return err
		}
		if invOld.Id == uuid.Nil {
			return models.ErrNotFound
		}
		// data.Qty = invOld.Qty + data.QtyChange
		// //update data
		// myMap := structs.Map(data)
		// myMap["updated_by"] = Claims.UserID
		// delete(myMap, "QtyChange")
		// delete(myMap, "Description")
		invOld.Qty += data.QtyChange
		invOld.UpdatedBy = uuid.FromStringOrNil(Claims.UserID)
		err = u.repoInventory.Update(ctx, ID, invOld)
		if err != nil {
			return err
		}
	}
	//data outlets
	outlets, err := u.repoOutlet.GetDataBy(ctx, "id", data.OutletId.String())
	if err != nil {
		return err
	}
	//data sku
	product, err := u.repoProduct.GetDataBy(ctx, "id", data.ProductId.String())
	if err != nil {
		return err
	}

	auditLogs := &models.AddAuditLogs{
		AuditDate:   util.GetTimeNow(),
		Source:      "outlet-inventory",
		Username:    Claims.UserName,
		OutletId:    data.OutletId,
		OutletName:  outlets.OutletName,
		ProductId:   data.ProductId,
		SkuName:     product.SkuName,
		Qty:         data.Qty,
		QtyChange:   data.QtyChange + data.Qty,
		QtyDelta:    data.QtyChange,
		Description: data.Description,
	}

	//save audit log
	err = u.useAuditLogs.Create(ctx, Claims, auditLogs)
	if err != nil {
		return err
	}

	return nil
}

func (u *useInventory) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoInventory.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}

// PatchStock implements iinventory.Usecase
func (u *useInventory) PatchStock(ctx context.Context, Claims util.Claims, param models.InvPatchStockRequest) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var logger = logging.Logger{}

	invList, err := u.repoInventory.GetList(ctx,
		models.ParamList{
			InitSearch: fmt.Sprintf("outlet_id = '%s' and product_id = '%s'", param.OutletId, param.ProductId),
		},
	)
	if err != nil {
		logger.Error("error inventory/PatchStock list ", err)
		return err
	}
	if len(invList) == 0 {
		//if no data then insert with zero qty
		err = u.repoInventory.Create(ctx, &models.Inventory{
			AddInventory: models.AddInventory{
				OutletId:  param.OutletId,
				ProductId: param.ProductId,
				Qty:       0 - param.Qty,
			},
		})
		if err != nil {
			logger.Error("error inventory/PatchStock create inventory ", err)
			return err
		}
		// return models.ErrInventoryNotFound
	} else {
		inv := invList[0]

		// if inv.Qty+param.Qty < 0 {
		// 	return models.ErrQtyExceedStock
		// }

		inv.Qty += param.Qty

		err = u.repoInventory.Update(ctx, inv.Id, inv)
		if err != nil {
			logger.Error("error inventory/PatchStock update inventory ", err)
			return err
		}
	}

	return nil
}
