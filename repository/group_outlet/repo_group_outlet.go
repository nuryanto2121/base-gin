
package repogroupoutlet

import (
	"context"
	"fmt"

	igroupoutlet "app/interface/group_outlet"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"
	
	"gorm.io/gorm"
	uuid "github.com/satori/go.uuid"
)
	
type repoGroupOutlet struct {
	Conn *gorm.DB
}
	
func NewRepoGroupOutlet(Conn *gorm.DB) igroupoutlet.Repository {
	return &repoGroupOutlet{Conn}
}
	
func (db *repoGroupOutlet) GetDataBy(ctx context.Context, ID uuid.UUID) (result *models.GroupOutlet, err error) {
	var (
		logger          = logging.Logger{}
		mGroupOutlet = &models.GroupOutlet{}
	)
	query := db.Conn.Where("group_outlet_id = ? ", ID).WithContext(ctx).Find(mGroupOutlet)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mGroupOutlet, nil
}
	
func (db *repoGroupOutlet) GetList(ctx context.Context,queryparam models.ParamList) (result []*models.GroupOutlet, err error) {

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

func (db *repoGroupOutlet) Create(ctx context.Context,data *models.GroupOutlet) error {
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
func (db *repoGroupOutlet) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.GroupOutlet{}).Where("groupoutlet_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoGroupOutlet) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("group_outlet_id = ?", ID).Delete(&models.GroupOutlet{})
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoGroupOutlet) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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
		query = db.Conn.Model(&models.GroupOutlet{}).Where(sWhere, queryparam.Search).Count(&rest)
	} else {
		query = db.Conn.Model(&models.GroupOutlet{}).Where(sWhere).Count(&rest)
	}
	// end where
	
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}
	
	return rest, nil
}
		
	