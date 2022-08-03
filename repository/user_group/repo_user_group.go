package repousergroup

import (
	"context"
	"fmt"

	iusergroup "app/interface/user_group"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type repoUserGroup struct {
	Conn *gorm.DB
}

func NewRepoUserGroup(Conn *gorm.DB) iusergroup.Repository {
	return &repoUserGroup{Conn}
}

func (db *repoUserGroup) GetById(ctx context.Context, ID uuid.UUID) (result *models.UserGroup, err error) {
	var (
		logger     = logging.Logger{}
		mUserGroup = &models.UserGroup{}
	)
	query := db.Conn.Where("user_group_id = ? ", ID).WithContext(ctx).Find(mUserGroup)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mUserGroup, nil
}

func (db *repoUserGroup) GetDataBy(ctx context.Context, key, value string) (result *models.UserGroupDesc, err error) {
	var (
		logger     = logging.Logger{}
		mUserGroup = &models.UserGroupDesc{}
	)
	query := db.Conn.Where(fmt.Sprintf("%s = ?", key)).WithContext(ctx).Find(mUserGroup)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mUserGroup, nil
}
func (db *repoUserGroup) GetListByUser(ctx context.Context, key, value string) (result []*models.UserGroupDesc, err error) {
	var (
		logger     = logging.Logger{}
		mUserGroup = []*models.UserGroupDesc{}
	)
	// ll := db.Conn.O
	//query := db.Conn.Where(fmt.Sprintf("%s = ?", key)).WithContext(ctx).Find(mUserGroup)
	query := db.Conn.Raw(`SELECT 
		ug.user_id,	ug.group_id ,
		g.group_code,	g.description ,
		to_jsonb(array_agg(otl)) as outlets
	 FROM user_group ug inner join groups g
		on ug.group_id =g.id 	 
	inner join 
	(select
			o.id as outlet_id,
			o.outlet_name,
			go2.group_id,
			go2.user_id
		from
			outlets o
		inner join group_outlet go2
			on	o.id = go2.outlet_id 		
	)as otl on
		otl.group_id = ug.group_id
		and otl.user_id = ug.user_id	
	 WHERE ug.user_id = ? GROUP BY ug.user_id,ug.group_id ,g.group_code,g.description ORDER BY g.group_code asc
	`, value).Find(&mUserGroup)
	// query := db.Conn.WithContext(ctx).Table(`user_group ug`).Select(`
	// ug.user_id,
	// ug.group_id ,
	// g.group_code,
	// g.description ,
	// json_agg(otl.outlet_name) as outlets
	// `).Joins(`inner join groups g
	// on ug.group_id =g.id
	// `).Joins(`inner join
	// (
	// 	select
	// 		o.outlet_name ,
	// 		o.outlet_city,
	// 		go2.group_id,
	// 		go2.user_id
	// 	from outlets o
	// 	inner join group_outlet go2
	// 		on o.id = go2.outlet_id
	// )as otl on
	// 	otl.group_id = ug.group_id
	// 	and otl.user_id = ug.user_id
	// `).Group(`ug.user_id,ug.group_id ,g.group_code,g.description`).Order(`g.group_code asc`).Where(fmt.Sprintf("ug.%s = ?", key)).Find(&mUserGroup)
	err = query.Error
	if err != nil {
		logger.Error("GetListByUser ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mUserGroup, nil
}

func (db *repoUserGroup) GetList(ctx context.Context, queryparam models.ParamList) (result []*models.UserGroup, err error) {

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

func (db *repoUserGroup) Create(ctx context.Context, data *models.UserGroup) error {
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
func (db *repoUserGroup) Update(ctx context.Context, ID uuid.UUID, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.UserGroup{}).Where("usergroup_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoUserGroup) Delete(ctx context.Context, ID uuid.UUID) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("user_id = ?", ID).Delete(&models.UserGroup{})
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoUserGroup) Count(ctx context.Context, queryparam models.ParamList) (result int64, err error) {
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
		query = db.Conn.Model(&models.UserGroup{}).Where(sWhere, queryparam.Search).Count(&rest)
	} else {
		query = db.Conn.Model(&models.UserGroup{}).Where(sWhere).Count(&rest)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return rest, nil
}
