package repoinventory

import (
	"context"

	iinventory "app/interface/inventory"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoInventory struct {
	Conn *gorm.DB
}

func NewRepoInventory(Conn *gorm.DB) iinventory.Repository {
	return &repoInventory{Conn}
}

func (db *repoInventory) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Inventory, err error) {
	var (
		logger     = logging.Logger{}
		mInventory = &models.Inventory{}
	)
	query := db.Conn.Where("id = ? ", ID).WithContext(ctx).Find(mInventory)

	err = query.Error
	if err != nil {
		logger.Error("repo inventory GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mInventory, nil
}

func (db *repoInventory) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Inventory, err error) {

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
			sWhere += " and (lower() LIKE ?)"
		} else {
			sWhere += "(lower() LIKE ?)"
		}
		query = db.Conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	}

	err = query.Error
	if err != nil {
		logger.Error("repo inventory GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (db *repoInventory) Create(ctx context.Context, data *models.Inventory) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Create(data)

	err = query.Error
	if err != nil {
		logger.Error("repo inventory Create ", err)
		return err
	}
	return nil
}
func (db *repoInventory) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.Inventory{}).Where("id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo inventory Update ", err)
		return err
	}
	return nil
}

func (db *repoInventory) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("id = ?", ID).Delete(&models.Inventory{})

	err = query.Error
	if err != nil {
		logger.Error("repo inventory Delete ", err)
		return err
	}
	return nil
}

func (db *repoInventory) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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
			sWhere += " and (lower() LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower() LIKE ? )" //queryparam.Search
		}
		query = db.Conn.Model(&models.Inventory{}).Where(sWhere, queryparam.Search).Count(&rest)
	} else {
		query = db.Conn.Model(&models.Inventory{}).Where(sWhere).Count(&rest)
	}
	// end where

	err = query.Error
	if err != nil {
		logger.Error("repo inventory Count ", err)
		return 0, err
	}

	return rest, nil
}
