package repopaymentnotificationlogs

import (
	"context"
	"fmt"

	imidtransnotificationlog "app/interface/payment_notification_logs"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoMidtransNotificationLog struct {
	db db.DBGormDelegate
}

func NewRepoMidtransNotificationLog(Conn db.DBGormDelegate) imidtransnotificationlog.Repository {
	return &repoMidtransNotificationLog{Conn}
}

func (r *repoMidtransNotificationLog) GetDataBy(ctx context.Context, key, value string) (*models.MidtransNotificationLog, error) {
	var (
		logger                   = logging.Logger{}
		mMidtransNotificationLog = &models.MidtransNotificationLog{}
		conn                     = r.db.Get(ctx)
	)

	err := conn.Where(fmt.Sprintf("%s = ?", key), value).WithContext(ctx).Find(mMidtransNotificationLog).Error
	if err != nil {
		logger.Error("repo midtrans_notification_log GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mMidtransNotificationLog, nil
}

func (r *repoMidtransNotificationLog) GetList(ctx context.Context, queryparam models.ParamList) ([]*models.MidtransNotificationLog, error) {

	var (
		pageNum  = 0
		pageSize = setting.AppSetting.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		conn     = r.db.Get(ctx)
		result   = []*models.MidtransNotificationLog{}
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
			sWhere += " and (lower() LIKE ?)"
		} else {
			sWhere += "(lower() LIKE ?)"
		}
		err = conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	} else {
		err = conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result).Error
	}

	if err != nil {
		logger.Error("repo midtrans_notification_log GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoMidtransNotificationLog) Create(ctx context.Context, data *models.MidtransNotificationLog) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Create(data).Error
	if err != nil {
		logger.Error("repo midtrans_notification_log Create ", err)
		return err
	}
	return nil
}
func (r *repoMidtransNotificationLog) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Model(models.MidtransNotificationLog{}).Where("midtransnotificationlog_id = ?", ID).Updates(data).Error
	if err != nil {
		logger.Error("repo midtrans_notification_log Update ", err)
		return err
	}
	return nil
}

func (r *repoMidtransNotificationLog) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
	)

	err := conn.Where("midtrans_notification_log_id = ?", ID).Delete(&models.MidtransNotificationLog{}).Error
	if err != nil {
		logger.Error("repo midtrans_notification_log Delete ", err)
		return err
	}
	return nil
}

func (r *repoMidtransNotificationLog) Count(ctx context.Context, queryparam models.ParamList) (int64, error) {
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
			sWhere += " and (lower() LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower() LIKE ? )" //queryparam.Search
		}
		err = conn.Model(&models.MidtransNotificationLog{}).Where(sWhere, queryparam.Search).Count(&rest).Error
	} else {
		err = conn.Model(&models.MidtransNotificationLog{}).Where(sWhere).Count(&rest).Error
	}
	// end where

	if err != nil {
		logger.Error("repo midtrans_notification_log Count ", err)
		return 0, err
	}

	return rest, nil
}
