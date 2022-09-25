package repotransactiondetail

import (
	"context"
	"fmt"

	itransactiondetail "app/interface/transaction_detail"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoTransactionDetail struct {
	db db.DBGormDelegate
}

func NewRepoTransactionDetail(Conn db.DBGormDelegate) itransactiondetail.Repository {
	return &repoTransactionDetail{Conn}
}

func (r *repoTransactionDetail) GetDataBy(ctx context.Context, key, value string) (*models.TransactionDetail, error) {
	var (
		logger             = logging.Logger{}
		mTransactionDetail = &models.TransactionDetail{}
		conn               = r.db.Get(ctx)
	)

	err := conn.Where(fmt.Sprintf("%s = ?", key), value).WithContext(ctx).Find(mTransactionDetail).Error
	if err != nil {
		logger.Error("repo transaction_detail GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mTransactionDetail, nil
}

func (r *repoTransactionDetail) GetList(ctx context.Context, queryparam models.ParamList) ([]*models.TransactionDetailRaw, error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		conn     = r.db.Get(ctx)
		result   = []*models.TransactionDetailRaw{}
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
			sWhere += " and (lower() LIKE ?)"
		} else {
			sWhere += "(lower() LIKE ?)"
		}
		err = conn.Model([]*models.TransactionDetail{}).
			Select(`transaction_detail.*,sku_management.sku_name`).
			Joins(`join sku_management sku_management on sku_management.id = transaction_detail.product_id`).
			Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
		// Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error

		// err = conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		// err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
		err = conn.Model([]*models.TransactionDetail{}).
			Select(`transaction_detail.*,sku_management.sku_name`).
			Joins(`join sku_management sku_management on sku_management.id = transaction_detail.product_id`).
			Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	if err != nil {
		logger.Error("repo transaction_detail GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoTransactionDetail) Create(ctx context.Context, data *models.TransactionDetail) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Create(data).Error
	if err != nil {
		logger.Error("repo transaction_detail Create ", err)
		return err
	}
	return nil
}
func (r *repoTransactionDetail) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Model(models.TransactionDetail{}).Where("id = ?", ID).Updates(data).Error
	if err != nil {
		logger.Error("repo transaction_detail Update ", err)
		return err
	}
	return nil
}

func (r *repoTransactionDetail) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Where("transaction_id = ?", ID).Delete(&models.TransactionDetail{}).Error
	if err != nil {
		logger.Error("repo transaction_detail Delete ", err)
		return err
	}
	return nil
}

func (r *repoTransactionDetail) Count(ctx context.Context, queryparam models.ParamList) (int64, error) {
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
			sWhere += " and (lower() LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower() LIKE ? )" //queryparam.Search
		}
		err = conn.Model(&models.TransactionDetail{}).Where(sWhere, queryparam.Search).Count(&rest).Error
	} else {
		err = conn.Model(&models.TransactionDetail{}).Where(sWhere).Count(&rest).Error
	}
	// end where

	if err != nil {
		logger.Error("repo transaction_detail Count ", err)
		return 0, err
	}

	return rest, nil
}
