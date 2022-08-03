
package repooutletdetail

import (
	"context"
	"fmt"

	ioutletdetail "app/interface/outlet_detail"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"
	
	"gorm.io/gorm"
	uuid "github.com/satori/go.uuid"
)
	
type repoOutletDetail struct {
	Conn *gorm.DB
}
	
func NewRepoOutletDetail(Conn *gorm.DB) ioutletdetail.Repository {
	return &repoOutletDetail{Conn}
}
	
func (db *repoOutletDetail) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.OutletDetail, err error) {
	var (
		logger          = logging.Logger{}
		mOutletDetail = &models.OutletDetail{}
	)
	query := db.Conn.Where("outlet_detail_id = ? ", ID).WithContext(ctx).Find(mOutletDetail)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mOutletDetail, nil
}
	
func (db *repoOutletDetail) GetList(ctx context.Context,queryparam models.ParamList) (result []*models.OutletDetail, err error) {

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

	
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (db *repoOutletDetail) Create(ctx context.Context,data *models.OutletDetail) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Create(data)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoOutletDetail) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.OutletDetail{}).Where("outletdetail_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoOutletDetail) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("outlet_detail_id = ?", ID).Delete(&models.OutletDetail{})
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoOutletDetail) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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
		query = db.Conn.Model(&models.OutletDetail{}).Where(sWhere, queryparam.Search).Count(&rest)
	} else {
		query = db.Conn.Model(&models.OutletDetail{}).Where(sWhere).Count(&rest)
	}
	// end where
	
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}
	
	return rest, nil
}
		
	