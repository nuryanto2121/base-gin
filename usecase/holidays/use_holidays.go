package useholidays

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	iholidays "app/interface/holidays"
	"app/models"

	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type useHolidays struct {
	repoHolidays   iholidays.Repository
	contextTimeOut time.Duration
}

func NewHolidaysHolidays(a iholidays.Repository, timeout time.Duration) iholidays.Usecase {
	return &useHolidays{repoHolidays: a, contextTimeOut: timeout}
}

func (u *useHolidays) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Holidays, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoHolidays.GetDataBy(ctx, ID)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *useHolidays) GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}
	result.Data, err = u.repoHolidays.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoHolidays.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useHolidays) Create(ctx context.Context, data *models.HolidayForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = &models.Holidays{}
	err = mapstructure.Decode(data, &form)
	if err != nil {
		return err
	}

	err = u.repoHolidays.Create(ctx, form)
	if err != nil {
		return err
	}
	return nil

}
func (u *useHolidays) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = models.Holidays{}
	dataOld, err := u.repoHolidays.GetDataBy(ctx, ID)
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

	err = u.repoHolidays.Update(ctx, ID, form)
	if err != nil {
		return err
	}
	return nil
}
func (u *useHolidays) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoHolidays.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
