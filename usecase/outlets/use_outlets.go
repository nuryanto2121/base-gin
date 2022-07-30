package useoutlets

import (
	ioutletDetail "app/interface/outlet_detail"
	ioutlets "app/interface/outlets"
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

type useOutlets struct {
	repoOutlets      ioutlets.Repository
	repoOutletDetail ioutletDetail.Repository
	contextTimeOut   time.Duration
}

func NewUseOutlets(a ioutlets.Repository, b ioutletDetail.Repository, timeout time.Duration) ioutlets.Usecase {
	return &useOutlets{
		repoOutlets:      a,
		repoOutletDetail: b,
		contextTimeOut:   timeout,
	}
}

func (u *useOutlets) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.Outlets, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoOutlets.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useOutlets) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoOutlets.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoOutlets.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useOutlets) Create(ctx context.Context, Claims util.Claims, data *models.OutletForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mOutlets = models.Outlets{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mOutlets.AddOutlets)
	if err != nil {
		return err
	}
	//check outlet by name
	isExist, err := u.repoOutlets.GetDataBy(ctx, "outlet_name", data.OutletName)
	if err != nil {
		return err
	}

	if isExist != nil && isExist.Id != uuid.Nil {
		return models.ErrConflict
	}

	mOutlets.CreatedBy = uuid.FromStringOrNil(Claims.Id)
	mOutlets.UpdatedBy = uuid.FromStringOrNil(Claims.Id)

	err = u.repoOutlets.Create(ctx, &mOutlets)
	if err != nil {
		return err
	}

	//insert detail
	for _, val := range data.OutletDetail {
		var mOutletDetail = models.OutletDetail{}
		val.ProductId = mOutlets.Id

		err = mapstructure.Decode(val, &mOutletDetail.AddOutletDetail)
		if err != nil {
			return err
		}

		mOutletDetail.CreatedBy = uuid.FromStringOrNil(Claims.Id)
		mOutletDetail.UpdatedBy = uuid.FromStringOrNil(Claims.Id)

		err = u.repoOutletDetail.Create(ctx, &mOutletDetail)
		if err != nil {
			return err
		}

	}
	return nil

}

func (u *useOutlets) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddOutlets) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoOutlets.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useOutlets) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoOutlets.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
