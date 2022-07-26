package useskumanagement

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	iskumanagement "app/interface/sku_management"
	"app/models"

	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type useskumanagement struct {
	reposkumanagement iskumanagement.Repository
	contextTimeOut    time.Duration
}

func NewSkuManagement(a iskumanagement.Repository, timeout time.Duration) iskumanagement.Usecase {
	return &useskumanagement{reposkumanagement: a, contextTimeOut: timeout}
}

func (u *useskumanagement) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.SkuManagement, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.reposkumanagement.GetDataBy(ctx, ID)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *useskumanagement) GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}
	result.Data, err = u.reposkumanagement.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.reposkumanagement.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useskumanagement) Create(ctx context.Context, data *models.SkuManagement) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = &models.SkuManagement{}
	err = mapstructure.Decode(data, &form)
	if err != nil {
		return err
	}

	err = u.reposkumanagement.Create(ctx, form)
	if err != nil {
		return err
	}
	return nil

}
func (u *useskumanagement) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = models.SkuManagement{}
	dataOld, err := u.reposkumanagement.GetDataBy(ctx, ID)
	if err != nil {
		return err
	}

	if dataOld.Id == uuid.Nil {
		return models.ErrNotFound
	}

	err = mapstructure.Decode(data, &form)
	if err != nil {
		return err
	}

	err = u.reposkumanagement.Update(ctx, ID, form)
	if err != nil {
		return err
	}
	return nil
}
func (u *useskumanagement) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.reposkumanagement.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
