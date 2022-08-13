package repotermandconditional

import (
	"context"
	"fmt"

	itermandconditional "app/interface/term_and_conditional"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoTermAndConditional struct {
	Conn *gorm.DB
}

func NewRepoTermAndConditioinal(Conn *gorm.DB) itermandconditional.Repository {
	return &repoTermAndConditional{Conn}
}

func (db *repoTermAndConditional) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.TermAndConditionalForm, err error) {
	var sysTermAndConditional = &models.TermAndConditionalForm{}
	query := db.Conn.WithContext(ctx).Where("id = ? ", ID).Find(sysTermAndConditional)
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return sysTermAndConditional, nil
}

func (db *repoTermAndConditional) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.TermAndConditionalForm, err error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		// logger   = logging.Logger{}
		orderBy = queryparam.SortField
	)
	// pagination
	if queryparam.Page > 0 {
		pageNum = (queryparam.Page - 1) * queryparam.PerPage
	}
	if queryparam.PerPage > 0 {
		pageSize = queryparam.PerPage
	}
	//end pagination

	// Order
	if queryparam.SortField != "" {
		orderBy = queryparam.SortField
	}
	//end Order by

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and " + queryparam.Search
		} else {
			sWhere += queryparam.Search
		}
	}

	// end where
	if pageNum >= 0 && pageSize > 0 {
		query := db.Conn.WithContext(ctx).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
		err = query.Error
	} else {
		query := db.Conn.WithContext(ctx).Where(sWhere).Order(orderBy).Find(&result)
		err = query.Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (db *repoTermAndConditional) Create(ctx context.Context, data *models.TermAndConditionalForm) (err error) {
	query := db.Conn.WithContext(ctx).Create(data)
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoTermAndConditional) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {

	query := db.Conn.WithContext(ctx).Model(models.TermAndConditionalForm{}).Where("id = ?", ID).Updates(data)
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

// func (db *repoHolidays) Delete(ctx context.Context, ID uuid.UUID) (err error) {

// 	query := db.Conn.WithContext(ctx).Where("id = ?", ID).Delete(&models.Holidays{})
// 	err = query.Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
func (db *repoTermAndConditional) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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

	query := db.Conn.WithContext(ctx).Model(&models.TermAndConditionalForm{}).Where(sWhere).Count(&result)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
