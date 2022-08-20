package repooutlets

import (
	"context"
	"fmt"

	ioutlets "app/interface/outlets"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoOutlets struct {
	db db.DBGormDelegate
}

func NewRepoOutlets(Conn db.DBGormDelegate) ioutlets.Repository {
	return &repoOutlets{Conn}
}

func (r *repoOutlets) GetDataBy(ctx context.Context, key, value string) (result *models.Outlets, err error) {
	var (
		logger   = logging.Logger{}
		mOutlets = &models.Outlets{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where(fmt.Sprintf("%s = ?", key), value).First(mOutlets) //(mOutlets)
	// logger.Query(fmt.Sprintf("%#v", query.Statement.Quote("")))
	err = query.Error
	if err != nil {
		logger.Error("repo outlet GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return mOutlets, models.ErrNotFound
		}
		return nil, err
	}
	return mOutlets, nil
}

func (r *repoOutlets) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.OutletList, err error) {

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
			sWhere += " and (lower(outlet_name) LIKE ?)"
		} else {
			sWhere += "(lower(outlet_name) LIKE ?)"
		}

		err = conn.Table(`outlets o`).Select(`
		 o.id as outlet_id
		 ,sm.id as product_id 
		 ,i.id as inventory_id
		 ,o.outlet_name 
		 ,o.outlet_city 
		 ,sm.sku_name 
		 ,coalesce(i.qty,0) as qty
		 ,sm.price_weekday 
		 ,sm.price_weekend 
		 ,od.outlet_price_weekday 
		 ,od.outlet_price_weekend 
		 ,ro.user_id 
		 ,ro.role 
		`).Joins(`cross join sku_management sm`).Joins(`
		inner join role_outlet ro
		 	on o.id = ro.outlet_id
		`).Joins(`
		left join outlet_detail od 
		 	on od.outlet_id = o.id 
		 	and od.product_id =sm.id
		`).Joins(`
		left join inventory i
		 	on i.outlet_id = o.id
			and i.product_id = sm.id
		`).Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {

		err = conn.Table(`outlets o`).Select(`
		 o.id as outlet_id
		 ,sm.id as product_id 
		 ,i.id as inventory_id
		 ,o.outlet_name 
		 ,o.outlet_city 
		 ,sm.sku_name 
		 ,coalesce(i.qty,0) as qty
		 ,sm.price_weekday 
		 ,sm.price_weekend 
		 ,od.outlet_price_weekday 
		 ,od.outlet_price_weekend 
		 ,ro.user_id 
		 ,ro.role 
		`).Joins(`cross join sku_management sm`).Joins(`
		inner join role_outlet ro
		 	on o.id = ro.outlet_id
		`).Joins(`
		left join outlet_detail od 
		 	on od.outlet_id = o.id 
		 	and od.product_id =sm.id
		`).Joins(`
		left join inventory i
		 	on i.outlet_id = o.id
			and i.product_id = sm.id
		`).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	if err != nil {
		logger.Error("repo outlet getlist ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (r *repoOutlets) Create(ctx context.Context, data *models.Outlets) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Create(data)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Create ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoOutlets) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Model(models.Outlets{}).Where("id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (r *repoOutlets) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Where("id = ?", ID).Delete(&models.Outlets{})

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (r *repoOutlets) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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

	sQuery := `
	select count(*) from 
	 (
		 select  o.id as outlet_id
		 ,sm.id as product_id 
		 ,i.id as inventory_id
		 ,o.outlet_name 
		 ,o.outlet_city 
		 ,sm.sku_name 
		 ,coalesce(i.qty,0) as qty
		 ,sm.price_weekday 
		 ,sm.price_weekend 
		 ,od.outlet_price_weekday 
		 ,od.outlet_price_weekend 
		 ,ro.user_id 
		 ,ro.role 
		 from outlets o 
		 cross join sku_management sm
		 inner join role_outlet ro
		 	on o.id = ro.outlet_id
		 left join outlet_detail od 
		 	on od.outlet_id = o.id 
		 	and od.product_id =sm.id 
		 left join inventory i
		 	on i.outlet_id = o.id
			and i.product_id = sm.id
	 ) outlet_list
	`

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and (lower(outlet_name) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(outlet_name) LIKE ? )" //queryparam.Search
		}
		sQuery += fmt.Sprintf(" WHERE %s", sWhere)
		err = conn.Raw(sQuery, queryparam.Search).Count(&rest).Error
	} else {
		sQuery += fmt.Sprintf(" WHERE %s", sWhere)
		err = conn.Raw(sQuery).Count(&rest).Error
	}
	// end where

	if err != nil {
		logger.Error("repo outlet Count ", err)
		return 0, models.ErrInternalServerError
	}

	return rest, nil
}
