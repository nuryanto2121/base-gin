package usesysuser

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	igroupoutlet "app/interface/group_outlet"
	iusers "app/interface/user"
	iusergroup "app/interface/user_group"
	"app/models"
	"app/pkg/logging"
	util "app/pkg/utils"

	uuid "github.com/satori/go.uuid"
)

type useSysUser struct {
	repoUser       iusers.Repository
	useUserGroup   iusergroup.Usecase
	useGroupOutlet igroupoutlet.Usecase
	contextTimeOut time.Duration
}

func NewUserSysUser(a iusers.Repository, b iusergroup.Usecase, c igroupoutlet.Usecase, timeout time.Duration) iusers.Usecase {
	return &useSysUser{
		repoUser:       a,
		useUserGroup:   b,
		useGroupOutlet: c,
		contextTimeOut: timeout}
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

	result, err = u.repoUser.GetById(ctx, ID)
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
func (u *useSysUser) CreateCms(ctx context.Context, req *models.AddUserCms) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var logger = logging.Logger{}

	if req.Password != req.ConfirmPassword {
		logger.Error(models.ErrWrongPasswordConfirm)
		return models.ErrWrongPasswordConfirm
	}

	//check username isexit
	isExist, err := u.repoUser.IsExist(ctx, "username", req.Username)
	if err != nil {
		logger.Error("IsExist ", err)
		return err
	}

	if isExist {
		return models.ErrAccountAlreadyExist
	}

	pass, _ := util.Hash(req.Password)
	dataUser := &models.Users{
		Username: req.Username,
		Name:     req.Username,
		Password: pass,
		IsActive: true,
	}
	err = u.repoUser.Create(ctx, dataUser)
	if err != nil {
		logger.Error("service user create ", err)
		return models.ErrInternalServerError
	}

	//save to user group
	dtCl := util.Claims{
		UserID: dataUser.Id.String(),
	}
	for _, val := range req.Groups {
		var dataGroupUser = models.AddUserGroup{
			UserId:  dataUser.Id,
			GroupId: val.GroupId,
		}
		//
		err := u.useUserGroup.Create(ctx, dtCl, &dataGroupUser)
		if err != nil {
			logger.Error("service save user group ", err)
			return models.ErrInternalServerError
		}

		for _, valOutlet := range val.OutletIds {
			err := u.useGroupOutlet.Create(ctx, dtCl, &models.AddGroupOutlet{
				GroupId:  val.GroupId,
				UserId:   dataUser.Id,
				OutletId: valOutlet.OutletId,
			})
			if err != nil {
				logger.Error("service save user group Create", err)
				return models.ErrInternalServerError
			}
		}
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
