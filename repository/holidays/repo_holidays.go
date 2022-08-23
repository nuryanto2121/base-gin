package repoholidays

import (
	"context"
	"fmt"

	iholidays "app/interface/holidays"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoHolidays struct {
	db db.DBGormDelegate
}

func NewRepoHolidays(Conn db.DBGormDelegate) iholidays.Repository {
	return &repoHolidays{Conn}
}

func (r *repoHolidays) GetDataBy(ctx context.Context, key, value string) (result *models.Holidays, err error) {
	var (
		sysHoliday = &models.Holidays{}
		logger     = logging.Logger{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where(fmt.Sprintf("%s = ?", key), value).First(sysHoliday)
	err = query.Error
	if err != nil {
		logger.Error("repo holiday GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return sysHoliday, models.ErrNotFound
		}
		return nil, err
	}
	return sysHoliday, nil
}

func (r *repoHolidays) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Holidays, err error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		// query    *gorm.DB
		conn = r.db.Get(ctx)
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
			sWhere += " and (lower(description) LIKE ?) "
		} else {
			sWhere += "(lower(description) LIKE ?)"
		}
		err = conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	// err = query.Error
	if err != nil {
		logger.Error("repo holiday GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (r *repoHolidays) Create(ctx context.Context, data *models.Holidays) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Create(data)
	err = query.Error
	if err != nil {
		logger.Error("repo holiday Create ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoHolidays) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Model(models.Holidays{}).Where("id = ?", ID).Updates(data)
	err = query.Error
	if err != nil {
		logger.Error("repo holiday Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoHolidays) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Where("id = ?", ID).Delete(&models.Holidays{})
	err = query.Error
	if err != nil {
		logger.Error("repo holiday Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoHolidays) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
		// query  *gorm.DB
		conn = r.db.Get(ctx)
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and (lower(description) LIKE ?) "
		} else {
			sWhere += "(lower(description) LIKE ?)"
		}
		err = conn.Model(&models.Holidays{}).Where(sWhere, queryparam.Search).Count(&result).Error
	} else {
		err = conn.Model(&models.Holidays{}).Where(sWhere).Count(&result).Error
	}
	// end where

	if err != nil {
		logger.Error("repo holiday Count ", err)
		return 0, models.ErrInternalServerError
	}

	return result, nil
}
