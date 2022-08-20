package repotermandconditional

import (
	"context"
	"fmt"

	itermandconditional "app/interface/term_and_conditional"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"

	// "app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoTermAndConditional struct {
	db db.DBGormDelegate
}

func NewRepoTermAndConditioinal(Conn db.DBGormDelegate) itermandconditional.Repository {
	return &repoTermAndConditional{Conn}
}

func (r *repoTermAndConditional) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.TermAndConditional, err error) {
	var sysTermAndConditional = &models.TermAndConditional{}
	conn := r.db.Get(ctx)
	query := conn.Where("id = ? ", ID).Find(sysTermAndConditional)
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return sysTermAndConditional, nil
}

func (r *repoTermAndConditional) GetDataOne(ctx context.Context) (result *models.TermAndConditional, err error) {
	var sysTermAndConditional = &models.TermAndConditional{}
	conn := r.db.Get(ctx)
	query := conn.First(sysTermAndConditional)
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return sysTermAndConditional, models.ErrNotFound
		}
		return nil, err
	}
	return sysTermAndConditional, nil
}

func (r *repoTermAndConditional) Create(ctx context.Context, data *models.TermAndConditional) (err error) {
	conn := r.db.Get(ctx)
	query := conn.Create(data)
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (r *repoTermAndConditional) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {

	conn := r.db.Get(ctx)
	query := conn.Model(models.TermAndConditional{}).Where("id = ?", ID).Updates(data)
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

// func (r *repoHolidays) Delete(ctx context.Context, ID uuid.UUID) (err error) {

// 	conn := r.db.Get(ctx)
//query := conn.Where("id = ?", ID).Delete(&models.Holidays{})
// 	err = query.Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
func (r *repoTermAndConditional) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and " + queryparam.Search
		}
	}
	// end where

	conn := r.db.Get(ctx)
	query := conn.Model(&models.TermAndConditional{}).Where(sWhere).Count(&result)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string

	if err != nil {
		return 0, err
	}

	return result, nil
}
