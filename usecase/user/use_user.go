package usesysuser

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	iusers "app/interface/user"
	"app/models"
	uuid "github.com/satori/go.uuid"
)

type useSysUser struct {
	repoUser       iusers.Repository
	contextTimeOut time.Duration
}

func NewUserSysUser(a iusers.Repository, timeout time.Duration) iusers.Usecase {
	return &useSysUser{repoUser: a, contextTimeOut: timeout}
}

func (u *useSysUser) GetByEmailSaUser(ctx context.Context, email string) (result *models.Users, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoUser.GetByAccount(ctx, email)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *useSysUser) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Users, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoUser.GetDataBy(ctx, ID)
	if err != nil {
		return result, err
	}
	result.Password = ""
	return result, nil
}
func (u *useSysUser) GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}
	result.Data, err = u.repoUser.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoUser.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useSysUser) Create(ctx context.Context, data *models.Users) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUser.Create(ctx, data)
	if err != nil {
		return err
	}
	return nil

}
func (u *useSysUser) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// var form = models.AddUser{}
	// err = mapstructure.Decode(data, &form)
	// if err != nil {
	// 	return err
	// 	// return appE.ResponseError(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)

	// }
	// err = u.repoUser.Update(ctx, ID, form)
	return nil
}
func (u *useSysUser) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUser.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
