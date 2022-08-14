package usesysuser

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	iroleoutlet "app/interface/role_outlet"
	iusers "app/interface/user"
	iuserrole "app/interface/user_role"
	"app/models"
	"app/pkg/logging"
	util "app/pkg/utils"

	uuid "github.com/satori/go.uuid"
)

type useSysUser struct {
	repoUser       iusers.Repository
	repoUserRole   iuserrole.Repository
	useRoleOutlet  iroleoutlet.Usecase
	contextTimeOut time.Duration
}

func NewUserSysUser(a iusers.Repository, b iuserrole.Repository, c iroleoutlet.Usecase, timeout time.Duration) iusers.Usecase {
	return &useSysUser{
		repoUser:       a,
		repoUserRole:   b,
		useRoleOutlet:  c,
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

func (u *useSysUser) GetDataBy(ctx context.Context, ID uuid.UUID) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	dataUser, err := u.repoUser.GetById(ctx, ID)
	if err != nil {
		return result, err
	}
	userList := []*models.ListUserCms{&models.ListUserCms{
		UserId:   dataUser.Id,
		Username: dataUser.Username,
	}}
	permission, err := genResponseList(u, ctx, userList)
	if err != nil {
		return result, err
	}
	if len(permission) <= 0 {
		return nil, models.ErrAccountNotFound
	}
	// }
	return permission[0], nil
}
func (u *useSysUser) GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}
	userList, err := u.repoUser.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Data, err = genResponseList(u, ctx, userList)
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
	for _, val := range req.Roles {
		var dataRoleUser = models.UserRole{
			AddUserRole: models.AddUserRole{
				UserId: dataUser.Id,
				Role:   val.Role,
			},
			Model: models.Model{
				CreatedBy: dataUser.Id,
				UpdatedBy: dataUser.Id,
			},
		}
		//
		err := u.repoUserRole.Create(ctx, &dataRoleUser)
		if err != nil {
			logger.Error("service save user role ", err)
			return models.ErrInternalServerError
		}

		for _, valOutlet := range val.OutletIds {
			err := u.useRoleOutlet.Create(ctx, dtCl, &models.AddRoleOutlet{
				Role:     val.Role,
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

func genResponseList(u *useSysUser, ctx context.Context, userList []*models.ListUserCms) ([]*models.ResponseListUserCms, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var result = []*models.ResponseListUserCms{}

	for _, val := range userList {
		var (
			userCms = &models.ResponseListUserCms{
				UserId:   val.UserId,
				Username: val.Username,
			}
		)
		//get user group/role
		userRoleList, err := u.repoUserRole.GetListByUser(ctx, "user_id", val.UserId.String())
		if err != nil {
			return nil, err
		}
		userCms.RoleCode = userRoleList

		result = append(result, userCms)
	}

	return result, nil
}
