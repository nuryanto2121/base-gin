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

type useUserRole struct {
	repoUserRole   iusergroup.Repository
	contextTimeOut time.Duration
}

func NewUseUserRole(a iusergroup.Repository, timeout time.Duration) iusergroup.Usecase {
	return &useUserRole{repoUserRole: a, contextTimeOut: timeout}
}

func (u *useUserRole) GetDataBy(ctx context.Context, Claims util.Claims, key, value string) (result *models.UserRoleDesc, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoUserRole.GetDataBy(ctx, key, value)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useUserRole) GetById(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.UserRole, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoUserRole.GetById(ctx, ID)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useUserRole) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoUserRole.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoUserRole.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useUserRole) Create(ctx context.Context, Claims util.Claims, data *models.AddUserRole) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mUserRole = models.UserRole{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mUserRole.AddUserRole)
	if err != nil {
		return err
	}

	mUserRole.CreatedBy = uuid.FromStringOrNil(Claims.Id)
	mUserRole.UpdatedBy = uuid.FromStringOrNil(Claims.Id)

	err = u.repoUserRole.Create(ctx, &mUserRole)
	if err != nil {
		return err
	}
	return nil

}

func (u *useUserRole) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddUserRole) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoUserRole.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useUserRole) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUserRole.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
func (u *useUserRole) DeleteByUserId(ctx context.Context, Claims util.Claims, UserID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUserRole.Delete(ctx, UserID)
	if err != nil {
		return err
	}
	return nil
}
