package repoauditlogs

import (
	"context"
	"fmt"

	iauditlogs "app/interface/audit_logs"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoAuditLogs struct {
	db db.DBGormDelegate
}

func NewRepoAuditLogs(Conn db.DBGormDelegate) iauditlogs.Repository {
	return &repoAuditLogs{Conn}
}

func (r *repoAuditLogs) GetDataBy(ctx context.Context, key, value string) (*models.AuditLogs, error) {
	var (
		logger     = logging.Logger{}
		mAuditLogs = &models.AuditLogs{}
		conn       = r.db.Get(ctx)
	)

	err := conn.Where(fmt.Sprintf("%s = ?", key), value).WithContext(ctx).Find(mAuditLogs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		logger.Error("repo audit_logs GetDataBy ", err)
		return nil, err
	}
	return mAuditLogs, nil
}

func (r *repoAuditLogs) GetList(ctx context.Context, queryparam models.ParamList) ([]*models.AuditLogs, error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		conn     = r.db.Get(ctx)
		result   = []*models.AuditLogs{}
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
			sWhere += " and (lower(username) LIKE ?)"
		} else {
			sWhere += "(lower(username) LIKE ?)"
		}
		err = conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	if err != nil {
		logger.Error("repo audit_logs GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoAuditLogs) Create(ctx context.Context, data *models.AuditLogs) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Create(data).Error
	if err != nil {
		logger.Error("repo audit_logs Create ", err)
		return err
	}
	return nil
}
func (r *repoAuditLogs) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Model(models.AuditLogs{}).Where("auditlogs_id = ?", ID).Updates(data).Error
	if err != nil {
		logger.Error("repo audit_logs Update ", err)
		return err
	}
	return nil
}

func (r *repoAuditLogs) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Where("audit_logs_id = ?", ID).Delete(&models.AuditLogs{}).Error
	if err != nil {
		logger.Error("repo audit_logs Delete ", err)
		return err
	}
	return nil
}

func (r *repoAuditLogs) Count(ctx context.Context, queryparam models.ParamList) (int64, error) {
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
			sWhere += " and (lower(username) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(username) LIKE ? )" //queryparam.Search
		}
		err = conn.Model(&models.AuditLogs{}).Where(sWhere, queryparam.Search).Count(&rest).Error
	} else {
		err = conn.Model(&models.AuditLogs{}).Where(sWhere).Count(&rest).Error
	}
	// end where

	if err != nil {
		logger.Error("repo audit_logs Count ", err)
		return 0, err
	}

	return rest, nil
}
