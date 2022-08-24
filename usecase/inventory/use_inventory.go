package useinventory

import (
	iinventory "app/interface/inventory"
	"app/models"
	util "app/pkg/util"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/jinzhu/copier"

	uuid "github.com/satori/go.uuid"
)

type useInventory struct {
	repoInventory  iinventory.Repository
	contextTimeOut time.Duration
}

func NewUseInventory(a iinventory.Repository, timeout time.Duration) iinventory.Usecase {
	return &useInventory{repoInventory: a, contextTimeOut: timeout}
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
		data.Qty = invOld.Qty + data.QtyChange
		//update data
		myMap := structs.Map(data)
		myMap["updated_by"] = Claims.UserID
		delete(myMap, "QtyChange")
		fmt.Println(myMap)
		err = u.repoInventory.Update(ctx, ID, myMap)
		if err != nil {
			return err
		}
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
