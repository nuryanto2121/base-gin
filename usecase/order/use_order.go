package useorder

import (
	iorder "app/interface/order"
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
	contextTimeOut time.Duration
}

func NewUseOrder(a iorder.Repository, timeout time.Duration) iorder.Usecase {
	return &useOrder{repoOrder: a, contextTimeOut: timeout}
}

func (u *useOrder) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.Order, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoOrder.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return result, err
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

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
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

	err = u.repoOrder.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
