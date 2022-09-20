package repotransaction

import (
	"context"
	"fmt"

	itransaction "app/interface/transaction"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoTransaction struct {
	db db.DBGormDelegate
}

func NewRepoTransaction(Conn db.DBGormDelegate) itransaction.Repository {
	return &repoTransaction{Conn}
}

func (r *repoTransaction) GetDataBy(ctx context.Context, key, value string) (*models.Transaction, error) {
	var (
		logger       = logging.Logger{}
		mTransaction = &models.Transaction{}
		conn         = r.db.Get(ctx)
	)

	err := conn.Where(fmt.Sprintf("%s = ?", key), value).WithContext(ctx).Find(mTransaction).Error
	if err != nil {
		logger.Error("repo transaction GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mTransaction, nil
}

// IsExist implements itransaction.Repository
func (r *repoTransaction) IsExist(ctx context.Context, sWhere string) (bool, error) {
	var (
		logger       = logging.Logger{}
		mTransaction = &models.Transaction{}
		conn         = r.db.Get(ctx)
	)

	err := conn.Where(sWhere).WithContext(ctx).Find(mTransaction).Error
	if err != nil {
		logger.Error("repo transaction GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return false, models.ErrNotFound
		}
		return false, err
	}
	return mTransaction.Id != uuid.Nil, nil
}
func (r *repoTransaction) GetList(ctx context.Context, queryparam models.ParamList) ([]*models.TransactionList, error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		conn     = r.db.Get(ctx)
		result   = []*models.TransactionList{}
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
			sWhere += " and (lower(ua.name) LIKE ?)"
		} else {
			sWhere += "(lower(ua.name) LIKE ?)"
		}
		// err = conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
		err = conn.Table(`"transaction" t`).Select(`
			ua."name" ,ua.phone_no ,ua.is_parent ,td.check_in,td.check_out , td.duration ,t.status_transaction ,t.status_payment
		`).
			Joins(`inner join transaction_detail td on t.id = td.transaction_id`).
			Joins(`inner join user_apps ua on td.customer_id = ua.id`).
			Where(sWhere, queryparam.Search).Offset(pageNum).
			Limit(pageSize).
			Order(orderBy).
			Find(&result).Error
	} else {
		// err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
		err = conn.Table(`"transaction" t`).Select(`
			ua."name" ,ua.phone_no ,ua.is_parent ,td.check_in,td.check_out , td.duration ,t.status_transaction ,t.status_payment
		`).
			Joins(`inner join transaction_detail td on t.id = td.transaction_id`).
			Joins(`inner join user_apps ua on td.customer_id = ua.id`).
			Where(sWhere).Offset(pageNum).
			Limit(pageSize).
			Order(orderBy).
			Find(&result).Error
	}

	if err != nil {
		logger.Error("repo transaction GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoTransaction) Create(ctx context.Context, data *models.Transaction) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Create(data).Error
	if err != nil {
		logger.Error("repo transaction Create ", err)
		return err
	}
	return nil
}
func (r *repoTransaction) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Model(models.Transaction{}).Where("id = ?", ID).Updates(data).Error
	if err != nil {
		logger.Error("repo transaction Update ", err)
		return err
	}
	return nil
}

func (r *repoTransaction) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Where("id = ?", ID).Delete(&models.Transaction{}).Error
	if err != nil {
		logger.Error("repo transaction Delete ", err)
		return err
	}
	return nil
}

func (r *repoTransaction) Count(ctx context.Context, queryparam models.ParamList) (int64, error) {
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
			sWhere += " and (lower(ua.name) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(ua.name) LIKE ? )" //queryparam.Search
		}
		// err = conn.Model(&models.Transaction{}).Where(sWhere, queryparam.Search).Count(&rest).Error
		err = conn.Table(`"transaction" t`).Select(`
			ua."name" ,ua.phone_no ,ua.is_parent ,td.check_in,td.check_out , td.duration ,t.status_transaction ,t.status_payment
		`).
			Joins(`inner join transaction_detail td on t.id = td.transaction_id`).
			Joins(`inner join user_apps ua on td.customer_id = ua.id`).
			Where(sWhere, queryparam.Search).
			Count(&rest).Error
	} else {
		// err = conn.Model(&models.Transaction{}).Where(sWhere).Count(&rest).Error
		err = conn.Table(`"transaction" t`).Select(`
			ua."name" ,ua.phone_no ,ua.is_parent ,td.check_in,td.check_out , td.duration ,t.status_transaction ,t.status_payment
		`).
			Joins(`inner join transaction_detail td on t.id = td.transaction_id`).
			Joins(`inner join user_apps ua on td.customer_id = ua.id`).
			Where(sWhere).
			Count(&rest).Error
	}
	// end where

	if err != nil {
		logger.Error("repo transaction Count ", err)
		return 0, err
	}

	return rest, nil
}

func (r *repoTransaction) GetListTicketUser(ctx context.Context, queryparam models.ParamList) ([]*models.TransactionResponse, error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		conn     = r.db.Get(ctx)
		result   = []*models.TransactionResponse{}
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
			sWhere += " and (lower(ua.name) LIKE ?)"
		} else {
			sWhere += "(lower(ua.name) LIKE ?)"
		}
		// err = conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
		err = conn.Table(`"transaction" t`).Select(`
			t.id, t.transaction_code , t.transaction_date
			,t.outlet_id ,o.outlet_name ,o.outlet_city 
			,t.total_ticket ,total_amount ,t.status_transaction 
			,case 
				when t.status_transaction = 2000001 then 'Booked' 
				when t.status_transaction = 2000002 then 'Active'
				when t.status_transaction = 2000003 then 'Selesai'
				when t.status_transaction = 2000004 then 'Draf'
				else ''
			end as status
		`).
			Joins(`join outlets o on t.outlet_id = o.id`).
			Where(sWhere, queryparam.Search).Offset(pageNum).
			Limit(pageSize).
			Order(orderBy).
			Find(&result).Error
	} else {
		// err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
		err = conn.Table(`"transaction" t`).Select(`
			t.id, t.transaction_code , t.transaction_date
			,t.outlet_id ,o.outlet_name ,o.outlet_city 
			,t.total_ticket ,total_amount ,t.status_transaction 
			,case 
				when t.status_transaction = 2000001 then 'Booked' 
				when t.status_transaction = 2000002 then 'Active'
				when t.status_transaction = 2000003 then 'Selesai'
				when t.status_transaction = 2000004 then 'Draf'
				else ''
			end as status
		`).
			Joins(`join outlets o on t.outlet_id = o.id`).
			Where(sWhere).Offset(pageNum).
			Limit(pageSize).
			Order(orderBy).
			Find(&result).Error
	}

	if err != nil {
		logger.Error("repo transaction GetList user ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoTransaction) CountUserList(ctx context.Context, queryparam models.ParamList) (int64, error) {
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
			sWhere += " and (lower(ua.name) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(ua.name) LIKE ? )" //queryparam.Search
		}

		err = conn.Table(`"transaction" t`).Select(`
			t.id, t.transaction_code , t.transaction_date
			,t.outlet_id ,o.outlet_name ,o.outlet_city 
			,t.total_ticket ,total_amount ,t.status_transaction 
			,case 
				when t.status_transaction = 2000001 then 'Booked' 
				when t.status_transaction = 2000002 then 'Active'
				else 'Selesai'
			end as status
		`).
			Joins(`join outlets o on t.outlet_id = o.id`).
			Where(sWhere, queryparam.Search).
			Count(&rest).Error
	} else {
		// err = conn.Model(&models.Transaction{}).Where(sWhere).Count(&rest).Error
		err = conn.Table(`"transaction" t`).Select(`
			t.id, t.transaction_code , t.transaction_date
			,t.outlet_id ,o.outlet_name ,o.outlet_city 
			,t.total_ticket ,total_amount ,t.status_transaction 
			,case 
				when t.status_transaction = 2000001 then 'Booked' 
				when t.status_transaction = 2000002 then 'Active'
				else 'Selesai'
			end as status
		`).
			Joins(`join outlets o on t.outlet_id = o.id`).
			Where(sWhere).
			Count(&rest).Error
	}
	// end where

	if err != nil {
		logger.Error("repo transaction Count ", err)
		return 0, err
	}

	return rest, nil
}
