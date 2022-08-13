package useholidays

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	iholidays "app/interface/holidays"
	"app/models"
	util "app/pkg/utils"

	"github.com/fatih/structs"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
)

type useHolidays struct {
	repoHolidays   iholidays.Repository
	contextTimeOut time.Duration
}

func NewHolidaysHolidays(a iholidays.Repository, timeout time.Duration) iholidays.Usecase {
	return &useHolidays{repoHolidays: a, contextTimeOut: timeout}
}

func (u *useHolidays) GetDataBy(ctx context.Context, claims util.Claims, ID uuid.UUID) (result *models.Holidays, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoHolidays.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *useHolidays) GetList(ctx context.Context, claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
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
func (u *useHolidays) Create(ctx context.Context, claims util.Claims, data *models.HolidayForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	//check is exist date
	dataOld, err := u.repoHolidays.GetDataBy(ctx, "DATE(holiday_date)", data.HolidayDate.Format("2006-01-02"))
	if err != nil {
		return err
	}
	if dataOld.Id != uuid.Nil {
		return models.ErrDataAlreadyExist
	}

	var form = &models.Holidays{}
	err = copier.Copy(&form, data)
	if err != nil {
		return err
	}

	form.CreatedBy = uuid.FromStringOrNil(claims.UserID)
	form.UpdatedBy = uuid.FromStringOrNil(claims.UserID)
	err = u.repoHolidays.Create(ctx, form)
	if err != nil {
		return err
	}
	return nil

}
func (u *useHolidays) Update(ctx context.Context, claims util.Claims, ID uuid.UUID, data *models.HolidayForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// var form = &models.Holidays{}
	dataOld, err := u.repoHolidays.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return err
	}

	if dataOld.Id == uuid.Nil {
		return models.ErrNotFound
	}

	myMap := structs.Map(data)
	myMap["updated_by"] = claims.UserID

	err = u.repoHolidays.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}
func (u *useHolidays) Delete(ctx context.Context, claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoHolidays.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
