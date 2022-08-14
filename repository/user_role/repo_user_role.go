package repouserrole

import (
	"context"
	"fmt"

	iuserrole "app/interface/user_role"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoUserRole struct {
	Conn *gorm.DB
}

func NewRepoUserRole(Conn *gorm.DB) iuserrole.Repository {
	return &repoUserRole{Conn}
}

func (db *repoUserRole) GetById(ctx context.Context, ID uuid.UUID) (result *models.UserRole, err error) {
	var (
		logger    = logging.Logger{}
		mUserRole = &models.UserRole{}
	)
	query := db.Conn.Where("user_role_id = ? ", ID).WithContext(ctx).Find(mUserRole)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return mUserRole, nil
}

func (db *repoUserRole) GetDataBy(ctx context.Context, key, value string) (result *models.UserRoleDesc, err error) {
	var (
		logger    = logging.Logger{}
		mUserRole = &models.UserRoleDesc{}
	)
	query := db.Conn.Where(fmt.Sprintf("%s = ?", key)).WithContext(ctx).Find(mUserRole)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return mUserRole, nil
}
func (db *repoUserRole) GetListByUser(ctx context.Context, key, value string) (result []*models.UserRoleDesc, err error) {
	var (
		logger    = logging.Logger{}
		mUserRole = []*models.UserRoleDesc{}
	)

	query := db.Conn.Raw(`
	SELECT ug.user_id,ug.role ,r.role_name 
		,to_jsonb(array_agg(otl)) as outlets
	FROM user_role ug 	inner join roles r
	 on r."role" = ug."role" 
	inner join 
	(select
			o.id as outlet_id,
			o.outlet_name,
			go2.role,
			go2.user_id
		from
			outlets o
		inner join role_outlet go2
			on	o.id = go2.outlet_id 		
	)as otl on
		otl.role = ug.role
		and otl.user_id = ug.user_id	
	 WHERE ug.user_id = ? GROUP BY ug.user_id,ug.role,r.role_name  ORDER BY ug.role asc
	`, value).Find(&mUserRole)

	err = query.Error
	if err != nil {
		logger.Error("GetListByUser ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return mUserRole, nil
}

func (db *repoUserRole) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.UserRole, err error) {

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
		logger.Error("repo user group GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return result, nil
}

func (db *repoUserRole) Create(ctx context.Context, data *models.UserRole) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Create(data)

	err = query.Error
	if err != nil {
		logger.Error("repo user group Create ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (db *repoUserRole) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.UserRole{}).Where("usergroup_id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo user group Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (db *repoUserRole) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("user_id = ?", ID).Delete(&models.UserRole{})

	err = query.Error
	if err != nil {
		logger.Error("repo user group Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (db *repoUserRole) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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
		query = db.Conn.Model(&models.UserRole{}).Where(sWhere, queryparam.Search).Count(&rest)
	} else {
		query = db.Conn.Model(&models.UserRole{}).Where(sWhere).Count(&rest)
	}
	// end where

	err = query.Error
	if err != nil {
		logger.Error("repo user group Count ", err)
		return 0, err
	}

	return rest, nil
}
