package repogroupoutlet

import (
	"context"
	"fmt"

	iroleoutlet "app/interface/role_outlet"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoRoleOutlet struct {
	db db.DBGormDelegate
}

func NewRepoRoleOutlet(Conn db.DBGormDelegate) iroleoutlet.Repository {
	return &repoRoleOutlet{Conn}
}

func (r *repoRoleOutlet) GetDataBy(ctx context.Context, key, value string) (result *models.RoleOutlet, err error) {
	var (
		logger      = logging.Logger{}
		mRoleOutlet = &models.RoleOutlet{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where(fmt.Sprintf("%s = ?", key), value).Find(mRoleOutlet)

	err = query.Error
	if err != nil {
		logger.Error("", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mRoleOutlet, nil
}

// func (r *repoRoleOutlet) GetListBy(ctx context.Context, key, value string) ([]*models.OutletLookUp, error) {
// 	var (
// 		logger      = logging.Logger{}
// 		mRoleOutlet = []*models.OutletLookUp{}
// 	)
// 	conn := r.db.Get(ctx)
// query := conn.Table(`from role_outlet a`).Where(fmt.Sprintf("%s = ?", key), value).Order(``).Find(&result)
// 	return mRoleOutlet, nil
// }

func (r *repoRoleOutlet) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.OutletLookUp, err error) {

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
			sWhere += " and (lower(o.outlet_name) LIKE ?)"
		} else {
			sWhere += "(lower(o.outlet_name) LIKE ?)"
		}
		err = conn.Table(`role_outlet a`).Select(`o.id as outlet_id,o.outlet_name,o.outlet_city`).
			Joins(`inner join outlets o on o.id =a.outlet_id`).Group(`o.id,o.outlet_name,o.outlet_city`).
			Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		err = conn.Table(`role_outlet a`).Select(`o.id as outlet_id,o.outlet_name,o.outlet_city`).
			Joins(`inner join outlets o on o.id =a.outlet_id`).Group(`o.id,o.outlet_name,o.outlet_city`).
			Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	if err != nil {
		logger.Error("repo outlet GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoRoleOutlet) Create(ctx context.Context, data *models.RoleOutlet) error {
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
func (r *repoRoleOutlet) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Model(models.RoleOutlet{}).Where("groupoutlet_id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (r *repoRoleOutlet) Delete(ctx context.Context, key, value string) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Where(fmt.Sprintf("%s = ?", key), value).Delete(&models.RoleOutlet{})

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (r *repoRoleOutlet) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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
			sWhere += " and (lower(o.outlet_name) LIKE ?)"
		} else {
			sWhere += "(lower(o.outlet_name) LIKE ?)"
		}
		err = conn.Table(`role_outlet a`).Select(`o.id as outlet_id,o.outlet_name`).
			Joins(`inner join outlets o on o.id =a.outlet_id`).Group(`o.id,o.outlet_name,o.outlet_city`).
			Where(sWhere, queryparam.Search).Count(&rest).Error
	} else {
		err = conn.Table(`role_outlet a`).Select(`o.id as outlet_id,o.outlet_name`).
			Joins(`inner join outlets o on o.id =a.outlet_id`).Group(`o.id,o.outlet_name,o.outlet_city`).
			Where(sWhere).Count(&rest).Error
	}
	// end where
	if err != nil {
		logger.Error("repo outlet Count ", err)
		return 0, models.ErrInternalServerError
	}

	return rest, nil
}
