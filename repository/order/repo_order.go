package repoorder

import (
	"context"
	"fmt"

	iorder "app/interface/order"
	"app/models"
	"app/pkg/logging"
	"app/pkg/postgres"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoOrder struct {
	db postgres.DBGormDelegate
}

func NewRepoOrder(Conn postgres.DBGormDelegate) iorder.Repository {
	return &repoOrder{Conn}
}

func (r *repoOrder) GetDataBy(ctx context.Context, key, value string) (result *models.Order, err error) {
	var (
		logger = logging.Logger{}
		mOrder = &models.Order{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where(fmt.Sprintf("%s = ?", key), value).WithContext(ctx).Find(mOrder)

	err = query.Error
	if err != nil {
		logger.Error("repo order GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mOrder, nil
}

func (r *repoOrder) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Order, err error) {

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
			sWhere += " and (lower(order_id) LIKE ?)"
		} else {
			sWhere += "(lower(order_id) LIKE ?)"
		}
		err = conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	// err = query.Error
	if err != nil {
		logger.Error("repo order GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoOrder) Create(ctx context.Context, data *models.Order) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Create(data)

	err = query.Error
	if err != nil {
		logger.Error("repo order Create ", err)
		return err
	}
	return nil
}
func (r *repoOrder) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Model(models.Order{}).Where("order_id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo order Update ", err)
		return err
	}
	return nil
}

func (r *repoOrder) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Where("order_id = ?", ID).Delete(&models.Order{})

	err = query.Error
	if err != nil {
		logger.Error("repo order Delete ", err)
		return err
	}
	return nil
}

func (r *repoOrder) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere         = ""
		logger         = logging.Logger{}
		rest   (int64) = 0
		conn           = r.db.Get(ctx)
	)

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and (lower(order_id) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(order_id) LIKE ? )" //queryparam.Search
		}
		err = conn.Model(&models.Order{}).Where(sWhere, queryparam.Search).Count(&rest).Error
	} else {
		err = conn.Model(&models.Order{}).Where(sWhere).Count(&rest).Error
	}
	// end where
	if err != nil {
		logger.Error("repo order Count ", err)
		return 0, err
	}

	return rest, nil
}
