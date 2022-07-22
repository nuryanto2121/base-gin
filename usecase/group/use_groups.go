package usegroup

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	igroup "app/interface/group"
	"app/models"

	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type useGroups struct {
	repoGroups     igroup.Repository
	contextTimeOut time.Duration
}

func NewGroups(a igroup.Repository, timeout time.Duration) igroup.Usecase {
	return &useGroups{repoGroups: a, contextTimeOut: timeout}
}

func (u *useGroups) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Groups, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoGroups.GetDataBy(ctx, ID)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *useGroups) GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}
	result.Data, err = u.repoGroups.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoGroups.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useGroups) Create(ctx context.Context, data *models.GroupForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = &models.GroupForm{}
	err = mapstructure.Decode(data, &form)
	if err != nil {
		return err
	}

	err = u.repoGroups.Create(ctx, form)
	if err != nil {
		return err
	}
	return nil

}
func (u *useGroups) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = models.GroupForm{}
	dataOld, err := u.repoGroups.GetDataBy(ctx, ID)
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

	err = u.repoGroups.Update(ctx, ID, form)
	if err != nil {
		return err
	}
	return nil
}
func (u *useGroups) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoGroups.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
