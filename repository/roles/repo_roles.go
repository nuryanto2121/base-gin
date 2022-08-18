package repogroups

import (
	"context"

	irole "app/interface/role"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoRoles struct {
	Conn *gorm.DB
}

func NewRepoRoles(Conn *gorm.DB) irole.Repository {
	return &repoRoles{Conn}
}

func (db *repoRoles) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Roles, err error) {
	var (
		sysRoles = &models.Roles{}
		logger   = logging.Logger{}
	)
	query := db.Conn.WithContext(ctx).Where("id = ? ", ID).First(sysRoles)
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

func (db *repoRoles) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Roles, err error) {

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
			sWhere += " and ((lower(role) LIKE ?) OR (lower(role_name) LIKE ?))"
		} else {
			sWhere += "((lower(role) LIKE ?) OR (lower(role_name) LIKE ?))"
		}
		query = db.Conn.WithContext(ctx).Where(sWhere, queryparam.Search, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.WithContext(ctx).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	}
	err = query.Error
	if err != nil {
		logger.Error("repo role getlist ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (db *repoRoles) Create(ctx context.Context, data *models.Roles) (err error) {
	var logger = logging.Logger{}
	query := db.Conn.WithContext(ctx).Create(data)
	err = query.Error
	if err != nil {
		logger.Error("repo role Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (db *repoRoles) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	var logger = logging.Logger{}
	query := db.Conn.WithContext(ctx).Model(models.Roles{}).Where("id = ?", ID).Updates(data)
	err = query.Error
	if err != nil {
		logger.Error("repo role Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (db *repoRoles) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	var logger = logging.Logger{}
	query := db.Conn.WithContext(ctx).Where("id = ?", ID).Delete(&models.Roles{})
	err = query.Error
	if err != nil {
		logger.Error("repo role Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (db *repoRoles) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
		query  *gorm.DB
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
		query = db.Conn.WithContext(ctx).Model(&models.Roles{}).Where(sWhere, queryparam.Search, queryparam.Search).Count(&result)
	} else {
		query = db.Conn.WithContext(ctx).Model(&models.Roles{}).Where(sWhere).Count(&result)
	}
	// end where

	// query := db.Conn.WithContext(ctx).Model(&models.Roles{}).Where(sWhere).Count(&result)

	err = query.Error
	if err != nil {
		logger.Error("repo role count ", err)
		return 0, models.ErrInternalServerError
	}

	return result, nil
}
