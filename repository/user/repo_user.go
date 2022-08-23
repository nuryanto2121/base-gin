package repousers

import (
	"context"
	"fmt"

	iusers "app/interface/user"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoSysUser struct {
	db db.DBGormDelegate
}

func NewRepoSysUser(Conn db.DBGormDelegate) iusers.Repository {
	return &repoSysUser{Conn}
}
func (r *repoSysUser) GetByAccount(ctx context.Context, Account string) (result *models.Users, err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Where("(username like ? OR phone_no = ?)", Account, Account).First(&result)
	err = query.Error
	if err != nil {
		logger.Error("repo users GetByAccount ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return result, err
}

func (r *repoSysUser) GetById(ctx context.Context, ID uuid.UUID) (result *models.UserCms, err error) {
	var (
		sysUser = &models.UserCms{}
		logger  = logging.Logger{}
	)
	conn := r.db.Get(ctx)
	query := conn.Table("users u").Select(`
				u.id as user_id, u.username ,u.name ,u.phone_no ,u.email ,r.role ,r.role_name
			`).Joins(`
			inner join user_role ur 
			on u.id = ur.user_id 
			`).Joins(`inner join roles r 
				on ur."role" =r."role" `).
		Where("u.id = ?", ID).Find(sysUser)
	// conn := r.db.Get(ctx)
	//query := conn.Where("id = ? ", ID).Find(sysUser)
	err = query.Error
	if err != nil {
		logger.Error("repo users GetById ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return sysUser, nil

}
func (r *repoSysUser) GetDataBy(ctx context.Context, key, value string) (*models.Users, error) {
	var (
		logger = logging.Logger{}
		result = &models.Users{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where(fmt.Sprintf("%s = ?", key), value).First(result)

	err := query.Error
	if err != nil {
		logger.Error("repo user GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (r *repoSysUser) IsExist(ctx context.Context, key, value string) (bool, error) {
	var (
		logger       = logging.Logger{}
		result int64 = 0
	)
	conn := r.db.Get(ctx)
	query := conn.Model(&models.Users{}).Where(fmt.Sprintf("%s = ?", key), value).Count(&result) //.Find(result)
	err := query.Error
	if err != nil {
		logger.Error("repo user Count ", err)
		if err == gorm.ErrRecordNotFound {
			return false, models.ErrNotFound
		}
		return false, err
	}
	return result > 0, nil
}

func (r *repoSysUser) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.UserCms, err error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		// query    *gorm.DB
		orderBy = queryparam.SortField
		conn    = r.db.Get(ctx)
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
			sWhere += " and ((lower(u.username) LIKE ?) OR (lower(u.name) LIKE ?))"
		} else {
			sWhere += "((lower(u.username) LIKE ?) OR (lower(u.name) LIKE ?))"
		}

		err = conn.Table("users u").Select(`
				u.id as user_id, u.username ,u.name ,u.phone_no ,u.email ,r.role ,r.role_name
			`).Joins(`
			inner join user_role ur 
			on u.id = ur.user_id 
			`).Joins(`inner join roles r 
				on ur."role" =r."role" `).
			Where(sWhere, queryparam.Search, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		err = conn.Table("users u").Select(`
				u.id as user_id, u.username ,u.name ,u.phone_no ,u.email ,r.role ,r.role_name
			`).Joins(`
			inner join user_role ur 
			on u.id = ur.user_id 
			`).Joins(`inner join roles r 
				on ur."role" =r."role" `).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	// err = query.Error

	if err != nil {
		logger.Error("repo users getlist ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (r *repoSysUser) Create(ctx context.Context, data *models.Users) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Create(data)
	err = query.Error
	if err != nil {
		logger.Error("repo user Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoSysUser) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Model(models.Users{}).Where("id = ?", ID).Updates(data)
	err = query.Error
	if err != nil {
		logger.Error("repo user Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoSysUser) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	var logger = logging.Logger{}
	conn := r.db.Get(ctx)
	query := conn.Where("id = ?", ID).Delete(&models.Users{})
	err = query.Error
	if err != nil {
		logger.Error("repo user Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoSysUser) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
		// query  *gorm.DB
		conn = r.db.Get(ctx)
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and ((lower(u.username) LIKE ?) OR (lower(u.name) LIKE ?)) "
		} else {
			sWhere += "((lower(u.username) LIKE ?) OR (lower(u.name) LIKE ?))"
		}
		err = conn.Table("users u").Select(`
		u.id as user_id, u.username ,u.name ,u.phone_no ,u.email ,r.role ,r.role_name
	`).Joins(`
	inner join user_role ur 
			on u.id = ur.user_id 
	`).Joins(`inner join roles r 
				on ur."role" =r."role" `).
			Where(sWhere, queryparam.Search, queryparam.Search).Count(&result).Error
	} else {
		err = conn.Table("users u").Select(`
		u.id as user_id, u.username ,u.name ,u.phone_no ,u.email ,r.role ,r.role_name
	`).Joins(`
	inner join user_role ur 
			on u.id = ur.user_id 
	`).Joins(`inner join roles r 
				on ur."role" =r."role" `).
			Where(sWhere).Count(&result).Error
	}
	// end where

	// err = query.Error
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return result, nil
}
