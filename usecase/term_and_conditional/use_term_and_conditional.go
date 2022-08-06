package usetermandconditional

import (
	"context"
	// "fmt"
	// "math"
	// "strings"
	"time"

	itermandconditional "app/interface/term_and_conditional"
	"app/models"

	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type useTermAndConditional struct {
	repoTermAndConditional itermandconditional.Repository
	contextTimeOut         time.Duration
}

func NewTermAndConditional(a itermandconditional.Repository, timeout time.Duration) itermandconditional.Usecase {
	return &useTermAndConditional{repoTermAndConditional: a, contextTimeOut: timeout}
}

func (u *useTermAndConditional) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.TermAndConditional, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoTermAndConditional.GetDataBy(ctx, ID)
	if err != nil {
		return result, err
	}
	return result, nil
}

// func (u *useHolidays) GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error) {
// 	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
// 	defer cancel()

// 	if queryparam.Search != "" {
// 		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
// 	}
// 	result.Data, err = u.repoHolidays.GetList(ctx, queryparam)
// 	if err != nil {
// 		return result, err
// 	}

// 	result.Total, err = u.repoHolidays.Count(ctx, queryparam)
// 	if err != nil {
// 		return result, err
// 	}

// 	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
// 	result.Page = queryparam.Page

// 	return result, nil
// }
func (u *useTermAndConditional) Create(ctx context.Context, data *models.TermAndConditional) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = &models.TermAndConditional{}
	err = mapstructure.Decode(data, &form)
	if err != nil {
		return err
	}

	err = u.repoTermAndConditional.Create(ctx, form)
	if err != nil {
		return err
	}
	return nil

}
func (u *useTermAndConditional) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var form = models.TermAndConditional{}
	dataOld, err := u.repoTermAndConditional.GetDataBy(ctx, ID)
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

	err = u.repoTermAndConditional.Update(ctx, ID, form)
	if err != nil {
		return err
	}
	return nil
}

// func (u *useHolidays) Delete(ctx context.Context, ID uuid.UUID) (err error) {
// 	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
// 	defer cancel()

// 	err = u.repoHolidays.Delete(ctx, ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
