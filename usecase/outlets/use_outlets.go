package useoutlets

import (
	ioutletDetail "app/interface/outlet_detail"
	ioutlets "app/interface/outlets"
	iroleoutlet "app/interface/role_outlet"
	"app/models"
	"app/pkg/logging"
	util "app/pkg/utils"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"

	uuid "github.com/satori/go.uuid"
)

type useOutlets struct {
	repoOutlets      ioutlets.Repository
	repoOutletDetail ioutletDetail.Repository
	repoRoleOutlet   iroleoutlet.Repository
	contextTimeOut   time.Duration
}

func NewUseOutlets(a ioutlets.Repository,
	b ioutletDetail.Repository,
	c iroleoutlet.Repository,
	timeout time.Duration) ioutlets.Usecase {
	return &useOutlets{
		repoOutlets:      a,
		repoOutletDetail: b,
		repoRoleOutlet:   c,
		contextTimeOut:   timeout,
	}
}

func (u *useOutlets) GetDataByRole(ctx context.Context, Claims util.Claims, role string) (result []*models.Outlets, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoOutlets.GetDataByRole(ctx, Claims.UserID, role) //GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *useOutlets) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		dataHeader = &models.Outlets{}
		result     = &models.OutletForm{}
		queryparam models.ParamList
	)

	dataHeader, err := u.repoOutlets.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(dataHeader.AddOutlets, &result)
	if err != nil {
		return nil, err
	}

	queryparam.InitSearch = fmt.Sprintf("user_id = '%s' and o.id = '%s'", Claims.UserID, ID)
	queryparam.Page = 1
	queryparam.PerPage = 10000

	detail, err := u.repoOutlets.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	rest := map[string]interface{}{
		"dataHeader": result,
		"dataDetail": detail,
	}

	return rest, nil
}

func (u *useOutlets) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += fmt.Sprintf(" AND user_id = '%s' ", Claims.UserID)
	} else {
		queryparam.InitSearch = fmt.Sprintf(" user_id = '%s' ", Claims.UserID)
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
		userID   = uuid.FromStringOrNil(Claims.UserID)
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mOutlets.AddOutlets)
	if err != nil {
		return err
	}
	//check outlet by name
	isExist, err := u.repoOutlets.GetDataBy(ctx, "outlet_name", data.OutletName)
	if err != nil && err != models.ErrNotFound {
		return models.ErrDataAlreadyExist
	}

	if isExist != nil && isExist.Id != uuid.Nil {
		return models.ErrConflict
	}

	mOutlets.CreatedBy = userID
	mOutlets.UpdatedBy = userID

	err = u.repoOutlets.Create(ctx, &mOutlets)
	if err != nil {
		return err
	}

	//insert detail
	for _, val := range data.OutletDetail {
		var mOutletDetail = models.OutletDetail{}
		val.OutletId = mOutlets.Id

		err = mapstructure.Decode(val, &mOutletDetail.AddOutletDetail)
		if err != nil {
			return err
		}

		mOutletDetail.CreatedBy = userID
		mOutletDetail.UpdatedBy = userID

		err = u.repoOutletDetail.Create(ctx, &mOutletDetail)
		if err != nil {
			return err
		}

	}
	go func() {
		//for root insert to role_outlet
		logger := logging.Logger{}
		cxts := context.Background()
		roleOutlet := &models.RoleOutlet{
			AddRoleOutlet: models.AddRoleOutlet{
				Role:     "root",
				OutletId: mOutlets.Id,
				UserId:   userID,
			},
		}
		err := u.repoRoleOutlet.Create(cxts, roleOutlet)
		if err != nil {
			logger.Error("error create role outlet ", err)
		}

	}()
	return nil

}

func (u *useOutlets) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.OutletForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		mOutlets = models.Outlets{}
		userID   = uuid.FromStringOrNil(Claims.UserID)
	)

	//check Id is exist
	dataUpdateHeader, err := u.repoOutlets.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return err
	}

	err = mapstructure.Decode(data, &dataUpdateHeader.AddOutlets)
	if err != nil {
		return err
	}

	//update header
	dataUpdateHeader.UpdatedBy = userID
	err = u.repoOutlets.Update(ctx, ID, dataUpdateHeader)
	if err != nil {
		return err
	}

	//delete then insert detail
	err = u.repoOutletDetail.Delete(ctx, ID)
	if err != nil {
		return err
	}
	//insert detail
	for _, val := range data.OutletDetail {
		var mOutletDetail = models.OutletDetail{}
		val.OutletId = mOutlets.Id

		err = mapstructure.Decode(val, &mOutletDetail.AddOutletDetail)
		if err != nil {
			return err
		}

		mOutletDetail.CreatedBy = userID
		mOutletDetail.UpdatedBy = userID

		err = u.repoOutletDetail.Create(ctx, &mOutletDetail)
		if err != nil {
			return err
		}

	}

	return nil
}

func (u *useOutlets) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	//detele outlet header
	err = u.repoOutlets.Delete(ctx, ID)
	if err != nil {
		return err
	}

	//detele outlet details
	err = u.repoOutletDetail.Delete(ctx, ID)
	if err != nil {
		return err
	}

	//detele outlet role outlet
	u.repoRoleOutlet.Delete(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}
