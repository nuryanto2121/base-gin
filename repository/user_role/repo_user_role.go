package repouserrole

import (
	"context"
	"fmt"

	iuserrole "app/interface/user_role"
	"app/models"
	"app/pkg/logging"
	"app/pkg/postgres"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoUserRole struct {
	db postgres.DBGormDelegate
}

func NewRepoUserRole(Conn postgres.DBGormDelegate) iuserrole.Repository {
	return &repoUserRole{Conn}
}

func (r *repoUserRole) GetById(ctx context.Context, ID uuid.UUID) (result *models.UserRole, err error) {
	var (
		logger    = logging.Logger{}
		mUserRole = &models.UserRole{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where("id = ? ", ID).WithContext(ctx).Find(mUserRole)
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

func (r *repoUserRole) GetDataBy(ctx context.Context, key, value string) (result *models.UserRoleDesc, err error) {
	var (
		logger    = logging.Logger{}
		mUserRole = &models.UserRoleDesc{}
	)
	conn := r.db.Get(ctx)
	query := conn.Table(`user_role`).Where(fmt.Sprintf("%s = ?", key), value).WithContext(ctx).Find(mUserRole)
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
func (r *repoUserRole) GetListByUser(ctx context.Context, key, value string) (result []*models.UserRoleDesc, err error) {
	var (
		logger    = logging.Logger{}
		mUserRole = []*models.UserRoleDesc{}
	)

	conn := r.db.Get(ctx)
	query := conn.Raw(`
	SELECT ug.user_id,ug.role ,r.role_name 
	FROM user_role ug 	inner join roles r
	 on r."role" = ug."role" 	
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

func (r *repoUserRole) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.UserRole, err error) {

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

	if err != nil {
		logger.Error("repo user group GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return result, nil
}

func (r *repoUserRole) Create(ctx context.Context, data *models.UserRole) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Create(data)

	err = query.Error
	if err != nil {
		logger.Error("repo user group Create ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoUserRole) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Model(models.UserRole{}).Where("user_id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo user group Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (r *repoUserRole) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Where("user_id = ?", ID).Delete(&models.UserRole{})

	err = query.Error
	if err != nil {
		logger.Error("repo user group Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (r *repoUserRole) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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
			sWhere += " and (lower() LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower() LIKE ? )" //queryparam.Search
		}
		err = conn.Model(&models.UserRole{}).Where(sWhere, queryparam.Search).Count(&rest).Error
	} else {
		err = conn.Model(&models.UserRole{}).Where(sWhere).Count(&rest).Error
	}
	// end where

	if err != nil {
		logger.Error("repo user group Count ", err)
		return 0, err
	}

	return rest, nil
}
