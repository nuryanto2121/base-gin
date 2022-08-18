package reposkumanagement

import (
	"context"
	"fmt"

	iskumanagement "app/interface/sku_management"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type reposkumanagement struct {
	Conn *gorm.DB
}

func NewRepoSkuManagement(Conn *gorm.DB) iskumanagement.Repository {
	return &reposkumanagement{Conn}
}

func (db *reposkumanagement) GetDataBy(ctx context.Context, key, value string) (result *models.SkuManagement, err error) {
	var (
		sysSkuManagement = &models.SkuManagement{}
		logger           = logging.Logger{}
	)
	query := db.Conn.WithContext(ctx).Where(fmt.Sprintf("%s = ?", key), value).First(sysSkuManagement)
	err = query.Error
	if err != nil {
		logger.Error("repo sku management GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return sysSkuManagement, models.ErrNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return sysSkuManagement, nil
}

func (db *reposkumanagement) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.SkuManagement, err error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		query    *gorm.DB
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
			sWhere += " and (lower(sku_name) LIKE ?)"
		} else {
			sWhere += "(lower(sku_name) LIKE ?)"
		}
		query = db.Conn.WithContext(ctx).Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.WithContext(ctx).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	}

	err = query.Error

	if err != nil {
		logger.Error("repo sku management GetList ", err)

		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return result, nil
}
func (db *reposkumanagement) Create(ctx context.Context, data *models.SkuManagement) (err error) {
	var logger = logging.Logger{}
	query := db.Conn.WithContext(ctx).Create(data)
	err = query.Error
	if err != nil {
		logger.Error("repo sku management Create ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (db *reposkumanagement) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	var logger = logging.Logger{}
	query := db.Conn.WithContext(ctx).Model(models.SkuManagement{}).Where("id = ?", ID).Updates(data)
	err = query.Error
	if err != nil {
		logger.Error("repo sku management Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (db *reposkumanagement) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	var logger = logging.Logger{}
	query := db.Conn.WithContext(ctx).Where("id = ?", ID).Delete(&models.SkuManagement{})
	err = query.Error
	if err != nil {
		logger.Error("repo sku management Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (db *reposkumanagement) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
		query  *gorm.DB
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and (lower(sku_name) LIKE ?)"
		} else {
			sWhere += "(lower(sku_name) LIKE ?)"
		}
		query = db.Conn.WithContext(ctx).Model(models.SkuManagement{}).Where(sWhere, queryparam.Search).Count(&result)
	} else {
		query = db.Conn.WithContext(ctx).Model(models.SkuManagement{}).Where(sWhere).Count(&result)
	}
	// end where

	err = query.Error
	if err != nil {
		logger.Error("repo sku management Count ", err)
		return 0, models.ErrInternalServerError
	}

	return result, nil
}
