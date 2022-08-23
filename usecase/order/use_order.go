package useorder

import (
	iorder "app/interface/order"
	ioutlets "app/interface/outlets"
	iskumanagement "app/interface/sku_management"
	"app/models"
	util "app/pkg/utils"
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
	contextTimeOut time.Duration
}

func NewUseOrder(a iorder.Repository, b ioutlets.Repository, c iskumanagement.Repository, timeout time.Duration) iorder.Usecase {
	return &useOrder{
		repoOrder:      a,
		repoOutlet:     b,
		repoProduct:    c,
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
