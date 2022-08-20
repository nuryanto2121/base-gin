package repoinventory

import (
	"context"

	iinventory "app/interface/inventory"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoInventory struct {
	db db.DBGormDelegate
}

func NewRepoInventory(Conn db.DBGormDelegate) iinventory.Repository {
	return &repoInventory{Conn}
}

func (r *repoInventory) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Inventory, err error) {
	var (
		logger     = logging.Logger{}
		mInventory = &models.Inventory{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where("id = ? ", ID).Find(mInventory)

	err = query.Error
	if err != nil {
		logger.Error("repo inventory GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return mInventory, models.ErrNotFound
		}
		return nil, err
	}
	return mInventory, nil
}

func (r *repoInventory) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Inventory, err error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		conn     = r.db.Get(ctx)
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
		err = conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	// err = query.Error
	if err != nil {
		logger.Error("repo inventory GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoInventory) Create(ctx context.Context, data *models.Inventory) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Create(data)

	err = query.Error
	if err != nil {
		logger.Error("repo inventory Create ", err)
		return err
	}
	return nil
}
func (r *repoInventory) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Model(models.Inventory{}).Where("id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo inventory Update ", err)
		return err
	}
	return nil
}

func (r *repoInventory) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Where("id = ?", ID).Delete(&models.Inventory{})

	err = query.Error
	if err != nil {
		logger.Error("repo inventory Delete ", err)
		return err
	}
	return nil
}

func (r *repoInventory) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
		// query  *gorm.DB
		rest (int64) = 0
		conn         = r.db.Get(ctx)
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
		err = conn.Model(&models.Inventory{}).Where(sWhere, queryparam.Search).Count(&rest).Error
	} else {
		err = conn.Model(&models.Inventory{}).Where(sWhere).Count(&rest).Error
	}
	// end where

	if err != nil {
		logger.Error("repo inventory Count ", err)
		return 0, err
	}

	return rest, nil
}
