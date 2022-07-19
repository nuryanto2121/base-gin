package repousers

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	iusers "gitlab.com/369-engineer/369backend/account/interface/user"
	"gitlab.com/369-engineer/369backend/account/models"
	"gitlab.com/369-engineer/369backend/account/pkg/logging"
	"gitlab.com/369-engineer/369backend/account/pkg/setting"
	"gorm.io/gorm"
)

type repoSysUser struct {
	Conn *gorm.DB
}

func NewRepoSysUser(Conn *gorm.DB) iusers.Repository {
	return &repoSysUser{Conn}
}
func (db *repoSysUser) GetByAccount(ctx context.Context, Account string) (result *models.Users, err error) {

	query := db.Conn.WithContext(ctx).Where("(email like ? OR phone_no = ?)", Account, Account).First(&result)
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return result, err
}

func (db *repoSysUser) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.Users, err error) {
	var sysUser = &models.Users{}
	query := db.Conn.WithContext(ctx).Where("id = ? ", ID).Find(sysUser)
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return sysUser, nil
}

func (db *repoSysUser) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.Users, err error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		// logger   = logging.Logger{}
		orderBy = queryparam.SortField
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
			sWhere += " and " + queryparam.Search
		} else {
			sWhere += queryparam.Search
		}
	}

	// end where
	if pageNum >= 0 && pageSize > 0 {
		query := db.Conn.WithContext(ctx).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
		err = query.Error
	} else {
		query := db.Conn.WithContext(ctx).Where(sWhere).Order(orderBy).Find(&result)
		err = query.Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (db *repoSysUser) Create(ctx context.Context, data *models.Users) (err error) {
	query := db.Conn.WithContext(ctx).Create(data)
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoSysUser) Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error) {

	query := db.Conn.WithContext(ctx).Model(models.Users{}).Where("id = ?", ID).Updates(data)
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoSysUser) Delete(ctx context.Context, ID uuid.UUID) (err error) {

	query := db.Conn.WithContext(ctx).Where("id = ?", ID).Delete(&models.Users{})
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoSysUser) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and " + queryparam.Search
		}
	}
	// end where

	query := db.Conn.WithContext(ctx).Model(&models.Users{}).Where(sWhere).Count(&result)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
