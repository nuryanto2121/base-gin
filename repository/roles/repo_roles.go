package repogroups

import (
	"context"

	irole "app/interface/role"
	"app/models"
	"app/pkg/logging"
	"app/pkg/postgres"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoRoles struct {
	db postgres.DBGormDelegate
}

func NewRepoRoles(Conn postgres.DBGormDelegate) irole.Repository {
	return &repoRoles{Conn}
}

func (r *repoRoles) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Roles, err error) {
	var (
		sysRoles = &models.Roles{}
		logger   = logging.Logger{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where("id = ? ", ID).First(sysRoles)
	err = query.Error
	if err != nil {
		logger.Error("repo role GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return sysRoles, nil
}

func (r *repoRoles) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Roles, err error) {

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
			sWhere += " and ((lower(role) LIKE ?) OR (lower(role_name) LIKE ?))"
		} else {
			sWhere += "((lower(role) LIKE ?) OR (lower(role_name) LIKE ?))"
		}
		err = conn.Where(sWhere, queryparam.Search, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	if err != nil {
		logger.Error("repo role getlist ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (r *repoRoles) Create(ctx context.Context, data *models.Roles) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Create(data)
	err = query.Error
	if err != nil {
		logger.Error("repo role Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoRoles) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Model(models.Roles{}).Where("id = ?", ID).Updates(data)
	err = query.Error
	if err != nil {
		logger.Error("repo role Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoRoles) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Where("id = ?", ID).Delete(&models.Roles{})
	err = query.Error
	if err != nil {
		logger.Error("repo role Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoRoles) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	// if queryparam.Search != "" {
	// 	if sWhere != "" {
	// 		sWhere += " and " + queryparam.Search
	// 	}
	// }
	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and ((lower(role) LIKE ?) OR (lower(role_name) LIKE ?))"
		} else {
			sWhere += "((lower(role) LIKE ?) OR (lower(role_name) LIKE ?))"
		}
		err = conn.Model(&models.Roles{}).Where(sWhere, queryparam.Search, queryparam.Search).Count(&result).Error
	} else {
		err = conn.Model(&models.Roles{}).Where(sWhere).Count(&result).Error
	}
	// end where

	if err != nil {
		logger.Error("repo role count ", err)
		return 0, models.ErrInternalServerError
	}

	return result, nil
}
