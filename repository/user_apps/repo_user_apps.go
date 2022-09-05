package repouserapps

import (
	"context"
	"fmt"

	iuserapps "app/interface/user_apps"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoUserApps struct {
	db db.DBGormDelegate
}

func NewRepoUserApps(Conn db.DBGormDelegate) iuserapps.Repository {
	return &repoUserApps{Conn}
}

func (r *repoUserApps) GetDataBy(ctx context.Context, key, value string) (*models.UserApps, error) {
	var (
		logger    = logging.Logger{}
		mUserApps = &models.UserApps{}
		conn      = r.db.Get(ctx)
	)

	err := conn.Where(fmt.Sprintf("%s = ?", key), value).WithContext(ctx).First(mUserApps).Error
	if err != nil {
		logger.Error("repo user_apps GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mUserApps, nil
}

func (r *repoUserApps) GetByAccount(ctx context.Context, Account string) (result *models.UserApps, err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Where("is_parent = true AND phone_no = ?", Account).First(&result)
	err = query.Error
	if err != nil {
		logger.Error("repo users GetByAccount ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return result, err
}

func (r *repoUserApps) GetList(ctx context.Context, queryparam models.ParamList) ([]*models.UserApps, error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		conn     = r.db.Get(ctx)
		result   = []*models.UserApps{}
		err      error
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
			sWhere += " and (lower(name) LIKE ?)"
		} else {
			sWhere += "(lower(name) LIKE ?)"
		}
		err = conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	if err != nil {
		logger.Error("repo user_apps GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoUserApps) Create(ctx context.Context, data *models.UserApps) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Create(data).Error
	if err != nil {
		logger.Error("repo user_apps Create ", err)
		return err
	}
	return nil
}
func (r *repoUserApps) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Model(models.UserApps{}).Where("userapps_id = ?", ID).Updates(data).Error
	if err != nil {
		logger.Error("repo user_apps Update ", err)
		return err
	}
	return nil
}

func (r *repoUserApps) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Where("user_apps_id = ?", ID).Delete(&models.UserApps{}).Error
	if err != nil {
		logger.Error("repo user_apps Delete ", err)
		return err
	}
	return nil
}

func (r *repoUserApps) Count(ctx context.Context, queryparam models.ParamList) (int64, error) {
	var (
		sWhere         = ""
		logger         = logging.Logger{}
		rest   (int64) = 0
		conn           = r.db.Get(ctx)
		err    error
	)

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and (lower(name) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(name) LIKE ? )" //queryparam.Search
		}
		err = conn.Model(&models.UserApps{}).Where(sWhere, queryparam.Search).Count(&rest).Error
	} else {
		err = conn.Model(&models.UserApps{}).Where(sWhere).Count(&rest).Error
	}
	// end where

	if err != nil {
		logger.Error("repo user_apps Count ", err)
		return 0, err
	}

	return rest, nil
}
