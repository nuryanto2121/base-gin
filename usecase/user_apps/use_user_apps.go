package useuserapps

import (
	iuserapps "app/interface/user_apps"
	"app/models"
	"app/pkg/util"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"

	uuid "github.com/satori/go.uuid"
)

type useUserApps struct {
	repoUserApps   iuserapps.Repository
	contextTimeOut time.Duration
}

func NewUseUserApps(a iuserapps.Repository, timeout time.Duration) iuserapps.Usecase {
	return &useUserApps{repoUserApps: a, contextTimeOut: timeout}
}

func (u *useUserApps) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.UserApps, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoUserApps.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useUserApps) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += fmt.Sprintf(" and is_parent = false and parent_id = '%s'", Claims.UserID)
	} else {
		queryparam.InitSearch += fmt.Sprintf(" is_parent = false and parent_id = '%s'", Claims.UserID)
	}

	result.Data, err = u.repoUserApps.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoUserApps.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useUserApps) Create(ctx context.Context, Claims util.Claims, data *models.AddUserApps) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mUserApps = models.UserApps{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mUserApps.AddUserApps)
	if err != nil {
		return err
	}

	mUserApps.CreatedBy = uuid.FromStringOrNil(Claims.Id)
	mUserApps.UpdatedBy = uuid.FromStringOrNil(Claims.Id)

	err = u.repoUserApps.Create(ctx, &mUserApps)
	if err != nil {
		return err
	}
	return nil

}

func (u *useUserApps) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddUserApps) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoUserApps.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useUserApps) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUserApps.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
