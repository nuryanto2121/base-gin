package repooutlets

import (
	"context"
	"fmt"

	ioutlets "app/interface/outlets"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoOutlets struct {
	Conn *gorm.DB
}

func NewRepoOutlets(Conn *gorm.DB) ioutlets.Repository {
	return &repoOutlets{Conn}
}

func (db *repoOutlets) GetDataByRole(ctx context.Context, ID, role string) (result []*models.Outlets, err error) {
	var (
		logger = logging.Logger{}
		// mOutlets = &models.Outlets{}
	)
	query := db.Conn.WithContext(ctx).Table(`outlets as outlets`).Select(`outlets.*`).Joins(`
	INNER JOIN role_outlet ro on outlets.id =ro.outlet_id
	`).Where(`ro.user_id = ? and ro.role = ?`, ID, role).Order(`outlets.outlet_name`).Find(&result)

	if err := query.Error; err != nil {
		logger.Error("repo outlet GetDataByRole ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (db *repoOutlets) GetDataBy(ctx context.Context, key, value string) (result *models.Outlets, err error) {
	var (
		logger   = logging.Logger{}
		mOutlets = &models.Outlets{}
	)
	query := db.Conn.Where(fmt.Sprintf("%s = ?", key), value).WithContext(ctx).First(mOutlets) //(mOutlets)
	// logger.Query(fmt.Sprintf("%#v", query.Statement.Quote("")))
	err = query.Error
	if err != nil {
		logger.Error("repo outlet GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mOutlets, nil
}

func (db *repoOutlets) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Outlets, err error) {

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
			sWhere += " and (lower(outlet_name) LIKE ?)"
		} else {
			sWhere += "(lower(outlet_name) LIKE ?)"
		}
		query = db.Conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	}

	err = query.Error
	if err != nil {
		logger.Error("repo outlet getlist ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (db *repoOutlets) Create(ctx context.Context, data *models.Outlets) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Create(data)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Create ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (db *repoOutlets) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.Outlets{}).Where("outlets_id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (db *repoOutlets) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("outlets_id = ?", ID).Delete(&models.Outlets{})

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (db *repoOutlets) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
		query  *gorm.DB
		rest   (int64) = 0
	)

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and (lower(outlet_name) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(outlet_name) LIKE ? )" //queryparam.Search
		}
		query = db.Conn.Model(&models.Outlets{}).Where(sWhere, queryparam.Search).Count(&rest)
	} else {
		query = db.Conn.Model(&models.Outlets{}).Where(sWhere).Count(&rest)
	}
	// end where

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Count ", err)
		return 0, models.ErrInternalServerError
	}

	return rest, nil
}
