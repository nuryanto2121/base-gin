package usetermandconditional

import (
	"context"
	// "fmt"
	// "math"
	// "strings"
	"time"

	itermandconditional "app/interface/term_and_conditional"
	"app/models"
	util "app/pkg/utils"

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

func (u *useTermAndConditional) GetDataOne(ctx context.Context, claims util.Claims) (result *models.TermAndConditional, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoTermAndConditional.GetDataOne(ctx)
	if err != nil && err != models.ErrNotFound {
		return result, err
	}
	if result.Id == uuid.Nil {
		return nil, nil
	}
	return result, nil
}

func (u *useTermAndConditional) Create(ctx context.Context, claims util.Claims, data *models.TermAndConditionalForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if data.Id == uuid.Nil {
		var form = &models.TermAndConditional{}
		err = mapstructure.Decode(data, &form)

		if err != nil {
			return err
		}

		form.CreatedBy = uuid.FromStringOrNil(claims.UserID)
		form.UpdatedBy = uuid.FromStringOrNil(claims.UserID)

		err = u.repoTermAndConditional.Create(ctx, form)
		if err != nil {
			return err
		}
	} else {
		dataOld, err := u.repoTermAndConditional.GetDataBy(ctx, data.Id)
		if err != nil {
			return err
		}

		if dataOld.Id == uuid.Nil {
			return models.ErrNotFound
		}

		TandCupdate := map[string]interface{}{
			"description": data.Description,
			"updated_by":  claims.UserID,
		}

		err = u.repoTermAndConditional.Update(ctx, data.Id, TandCupdate)
		if err != nil {
			return err
		}
	}

	return nil

}
func (u *useTermAndConditional) Update(ctx context.Context, claims util.Claims, ID uuid.UUID, data interface{}) (err error) {
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

	err = u.repoTermAndConditional.Update(ctx, ID, &form)
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
