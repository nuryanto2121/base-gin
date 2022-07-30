package useusergroup

import (
	iusergroup "app/interface/user_group"
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

type useUserGroup struct {
	repoUserGroup  iusergroup.Repository
	contextTimeOut time.Duration
}

func NewUseUserGroup(a iusergroup.Repository, timeout time.Duration) iusergroup.Usecase {
	return &useUserGroup{repoUserGroup: a, contextTimeOut: timeout}
}

func (u *useUserGroup) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.UserGroup, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoUserGroup.GetDataBy(ctx, ID)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useUserGroup) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoUserGroup.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoUserGroup.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useUserGroup) Create(ctx context.Context, Claims util.Claims, data *models.AddUserGroup) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mUserGroup = models.UserGroup{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mUserGroup.AddUserGroup)
	if err != nil {
		return err
	}

	mUserGroup.CreatedBy = uuid.FromStringOrNil(Claims.Id)
	mUserGroup.UpdatedBy = uuid.FromStringOrNil(Claims.Id)

	err = u.repoUserGroup.Create(ctx, &mUserGroup)
	if err != nil {
		return err
	}
	return nil

}

func (u *useUserGroup) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddUserGroup) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoUserGroup.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useUserGroup) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUserGroup.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
func (u *useUserGroup) DeleteByUserId(ctx context.Context, Claims util.Claims, UserID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUserGroup.Delete(ctx, UserID)
	if err != nil {
		return err
	}
	return nil
}
