package repooutletdetail

import (
	"context"

	ioutletdetail "app/interface/outlet_detail"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoOutletDetail struct {
	db db.DBGormDelegate
}

func NewRepoOutletDetail(Conn db.DBGormDelegate) ioutletdetail.Repository {
	return &repoOutletDetail{Conn}
}

func (r *repoOutletDetail) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.OutletDetail, err error) {
	var (
		logger        = logging.Logger{}
		mOutletDetail = &models.OutletDetail{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where("id = ? ", ID).Find(mOutletDetail)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet detail GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mOutletDetail, nil
}
func (r *repoOutletDetail) GetListBy(ctx context.Context, ID uuid.UUID) (result []*models.OutletDetail, err error) {
	var (
		logger        = logging.Logger{}
		mOutletDetail = []*models.OutletDetail{}
	)
	conn := r.db.Get(ctx)
	query := conn.Where("outlet_id = ? ", ID).Find(mOutletDetail)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet detail GetDataBy ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mOutletDetail, nil
}

func (r *repoOutletDetail) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.OutletDetail, err error) {

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
		logger.Error("repo outlet detail GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (r *repoOutletDetail) Create(ctx context.Context, data *models.OutletDetail) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Create(data)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet detail Create ", err)
		return models.ErrInternalServerError
	}
	return nil
}
func (r *repoOutletDetail) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Model(models.OutletDetail{}).Where("outletdetail_id = ?", ID).Updates(data)

	err = query.Error
	if err != nil {
		logger.Error("repo outlet detail Update ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (r *repoOutletDetail) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Where("outlet_id = ?", ID).Delete(&models.OutletDetail{})

	err = query.Error
	if err != nil {
		logger.Error("repo outlet detail Delete ", err)
		return models.ErrInternalServerError
	}
	return nil
}

func (r *repoOutletDetail) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}

		rest (int64) = 0
		conn         = r.db.Get(ctx)
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
		err = conn.Model(&models.OutletDetail{}).Where(sWhere, queryparam.Search).Count(&rest).Error
	} else {
		err = conn.Model(&models.OutletDetail{}).Where(sWhere).Count(&rest).Error
	}
	// end where
	if err != nil {
		logger.Error("repo outlet detail Count ", err)
		return 0, models.ErrInternalServerError
	}

	return rest, nil
}
