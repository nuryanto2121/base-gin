package useorder

import (
	iauditlogs "app/interface/audit_logs"
	iinventory "app/interface/inventory"
	iorder "app/interface/order"
	ioutlets "app/interface/outlets"
	iskumanagement "app/interface/sku_management"
	itrx "app/interface/trx"
	"app/models"
	util "app/pkg/util"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"

	uuid "github.com/satori/go.uuid"
)

type useOrder struct {
	repoOrder      iorder.Repository
	repoOutlet     ioutlets.Repository
	repoProduct    iskumanagement.Repository
	useAuditLogs   iauditlogs.Usecase
	useInventory   iinventory.Usecase
	repoTrx        itrx.Repository
	contextTimeOut time.Duration
}

func NewUseOrder(
	repoOrder iorder.Repository,
	repoOutlet ioutlets.Repository,
	repoProduct iskumanagement.Repository,
	useAuditLogs iauditlogs.Usecase,
	useInventory iinventory.Usecase,
	repoTrx itrx.Repository,
	timeout time.Duration,
) iorder.Usecase {
	return &useOrder{
		repoOrder:      repoOrder,
		repoOutlet:     repoOutlet,
		repoProduct:    repoProduct,
		useAuditLogs:   useAuditLogs,
		useInventory:   useInventory,
		repoTrx:        repoTrx,
		contextTimeOut: timeout}
}

func (u *useOrder) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	dataOrder, err := u.repoOrder.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return result, err
	}
	//getdata sku
	dataProduct, err := u.repoProduct.GetDataBy(ctx, "id", dataOrder.ProductId.String())
	if err != nil {
		return result, err
	}
	//get data  outlet
	dataOutlet, err := u.repoOutlet.GetDataBy(ctx, "id", dataOrder.OutletId.String())
	if err != nil {
		return result, err
	}

	result = map[string]interface{}{
		"id":           dataOrder.Id,
		"order_id":     dataOrder.OrderID,
		"order_date":   dataOrder.OrderDate,
		"outlet_id":    dataOrder.OutletId,
		"outlet_name":  dataOutlet.OutletName,
		"product_id":   dataOrder.ProductId,
		"sku_name":     dataProduct.SkuName,
		"start_number": dataOrder.StartNumber,
		"end_number":   dataOrder.EndNumber,
		"qty":          dataOrder.Qty,
		"status":       dataOrder.Status,
	}
	return result, nil
}

func (u *useOrder) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoOrder.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoOrder.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useOrder) Create(ctx context.Context, Claims util.Claims, data *models.AddOrder) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mOrder = models.Order{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mOrder.AddOrder)
	if err != nil {
		return err
	}

	//gen order id
	if data.OrderID == "" {
		t := &models.TmpCode{Prefix: "ORD"}
		mOrder.OrderID = util.GenCode(t)
	}

	//check duplicate order id
	dataOld, err := u.repoOrder.GetDataBy(ctx, "order_id", mOrder.OrderID)
	if err != nil && err != models.ErrNotFound {
		return err
	}
	if dataOld != nil {
		return models.ErrDataAlreadyExist
	}

	mOrder.CreatedBy = uuid.FromStringOrNil(Claims.Id)
	mOrder.UpdatedBy = uuid.FromStringOrNil(Claims.Id)

	err = u.repoOrder.Create(ctx, &mOrder)
	if err != nil {
		return err
	}
	return nil

}

func (u *useOrder) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddOrder) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	//check data exist
	_, err = u.repoOrder.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return err
	}

	myMap := structs.Map(data)
	myMap["updated_by"] = Claims.UserID
	delete(myMap, "OrderID")
	fmt.Println(myMap)
	err = u.repoOrder.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useOrder) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	//check data exist
	_, err = u.repoOrder.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return err
	}

	err = u.repoOrder.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateStatus implements iorder.Usecase
func (u *useOrder) UpdateStatus(c context.Context, Claims util.Claims, data *models.InventoryStatusForm) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeOut)
	defer cancel()

	var (
		ID              = data.ID
		qty       int64 = 0
		qtyChange int64 = 0
		qtyDelta  int64 = 0
	)
	//check data exist
	order, err := u.repoOrder.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return err
	}

	order.Status = data.Status
	myMap := structs.Map(order.AddOrder)
	myMap["updated_by"] = Claims.UserID
	// delete(myMap, "OrderID")
	fmt.Println(myMap)
	err = u.repoOrder.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}

	errTx := u.repoTrx.Run(ctx, func(trxCtx context.Context) error {
		//data outlets
		outlets, err := u.repoOutlet.GetDataBy(trxCtx, "id", order.OutletId.String())
		if err != nil {
			return err
		}
		//data sku
		product, err := u.repoProduct.GetDataBy(trxCtx, "id", order.ProductId.String())
		if err != nil {
			return err
		}
		//when approve add stock in inventory
		if data.Status == models.APPROVE {
			//get data inventory outlet
			invOutletList, err := u.useInventory.GetList(trxCtx, Claims, models.ParamList{
				Page:       1,
				PerPage:    1,
				InitSearch: fmt.Sprintf("outlet_id ='%s' and product_id = '%s' ", order.OutletId, order.ProductId),
			})
			if err != nil {
				return err
			}

			err = u.useInventory.PatchStock(trxCtx, Claims, models.InvPatchStockRequest{
				OutletId:  order.OutletId,
				ProductId: order.ProductId,
				Qty:       order.Qty,
			})
			if err != nil {
				return err
			}
			inv := invOutletList.Data.([]*models.Inventory)
			qty = inv[0].Qty //order.Qty
			qtyDelta = order.Qty
			qtyChange = order.Qty + inv[0].Qty
			data.Description = "approve order"
		} else if data.Status == models.REJECT {
			qty = order.Qty
		}

		auditLogs := &models.AddAuditLogs{
			AuditDate:   util.GetTimeNow(),
			OutletId:    order.OutletId,
			OutletName:  outlets.OutletName,
			ProductId:   order.ProductId,
			SkuName:     product.SkuName,
			Qty:         qty,
			QtyChange:   qtyChange,
			QtyDelta:    qtyDelta,
			Source:      "order",
			Description: data.Description,
			Username:    Claims.UserName,
		}

		err = u.useAuditLogs.Create(trxCtx, Claims, auditLogs)
		if err != nil {
			return err
		}
		return nil
	})

	if errTx != nil {
		return errTx
	}

	return nil
}
