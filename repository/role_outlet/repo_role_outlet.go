package repogroupoutlet

import (
	"context"
	"fmt"

	iroleoutlet "app/interface/role_outlet"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoRoleOutlet struct {
	Conn *gorm.DB
}

func NewRepoRoleOutlet(Conn *gorm.DB) iroleoutlet.Repository {
	return &repoRoleOutlet{Conn}
}

func (db *repoRoleOutlet) GetDataBy(ctx context.Context, key, value string) (result *models.RoleOutlet, err error) {
	var (
		logger      = logging.Logger{}
		mRoleOutlet = &models.RoleOutlet{}
	)
	query := db.Conn.Where(fmt.Sprintf("%s = ?", key), value).WithContext(ctx).Find(mRoleOutlet)

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

// func (db *repoRoleOutlet) GetListBy(ctx context.Context, key, value string) ([]*models.OutletLookUp, error) {
// 	var (
// 		logger      = logging.Logger{}
// 		mRoleOutlet = []*models.OutletLookUp{}
// 	)
// 	query := db.Conn.Table(`from role_outlet a`).Where(fmt.Sprintf("%s = ?", key), value).Order(``).Find(&result)
// 	return mRoleOutlet, nil
// }

func (db *repoRoleOutlet) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.OutletLookUp, err error) {

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
			sWhere += " and (lower(o.outlet_name) LIKE ?)"
		} else {
			sWhere += "(lower(o.outlet_name) LIKE ?)"
		}
		query = db.Conn.Table(`role_outlet a`).Select(`o.id as outlet_id,o.outlet_name,o.outlet_city`).
			Joins(`inner join outlets o on o.id =a.outlet_id`).Group(`o.id,o.outlet_name,o.outlet_city`).
			Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.Table(`role_outlet a`).Select(`o.id as outlet_id,o.outlet_name,o.outlet_city`).
			Joins(`inner join outlets o on o.id =a.outlet_id`).Group(`o.id,o.outlet_name,o.outlet_city`).
			Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	}

	err = query.Error

	if err != nil {
		logger.Error("repo outlet GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (db *repoRoleOutlet) Create(ctx context.Context, data *models.RoleOutlet) error {
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
func (db *repoRoleOutlet) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.RoleOutlet{}).Where("groupoutlet_id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (db *repoRoleOutlet) Delete(ctx context.Context, key, value string) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where(fmt.Sprintf("%s = ?", key), value).Delete(&models.RoleOutlet{})

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (db *repoRoleOutlet) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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
			sWhere += " and (lower(o.outlet_name) LIKE ?)"
		} else {
			sWhere += "(lower(o.outlet_name) LIKE ?)"
		}
		query = db.Conn.Table(`role_outlet a`).Select(`o.id as outlet_id,o.outlet_name`).
			Joins(`inner join outlets o on o.id =a.outlet_id`).Group(`o.id,o.outlet_name,o.outlet_city`).
			Where(sWhere, queryparam.Search).Count(&rest)
	} else {
		query = db.Conn.Table(`role_outlet a`).Select(`o.id as outlet_id,o.outlet_name`).
			Joins(`inner join outlets o on o.id =a.outlet_id`).Group(`o.id,o.outlet_name,o.outlet_city`).
			Where(sWhere).Count(&rest)
	}
	// end where

	err = query.Error
	if err != nil {
		logger.Error("repo outlet Count ", err)
		return 0, models.ErrInternalServerError
	}

	return rest, nil
}