package usegroup

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	irole "app/interface/role"
	"app/models"

	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
)

type useRoles struct {
	repoRoles      irole.Repository
	contextTimeOut time.Duration
}

func NewRoles(a irole.Repository, timeout time.Duration) irole.Usecase {
	return &useRoles{repoRoles: a, contextTimeOut: timeout}
}

func (u *useRoles) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Roles, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoRoles.GetDataBy(ctx, ID)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *useRoles) GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}
	result.Data, err = u.repoRoles.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoRoles.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useRoles) Create(ctx context.Context, data *models.RoleForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = &models.Roles{}
	err = copier.Copy(&form, data)
	if err != nil {
		return err
	}

	err = u.repoRoles.Create(ctx, form)
	if err != nil {
		return err
	}
	return nil

}
func (u *useRoles) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = models.Roles{}
	dataOld, err := u.repoRoles.GetDataBy(ctx, ID)
	if err != nil {
		return err
	}

	if dataOld.Id == uuid.Nil {
		return models.ErrNotFound
	}

	err = copier.Copy(&form, data)
	if err != nil {
		return err
	}

	err = u.repoRoles.Update(ctx, ID, form)
	if err != nil {
		return err
	}
	return nil
}
func (u *useRoles) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoRoles.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
