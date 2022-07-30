
package usegroupoutlet

import (
	"context"
	"fmt"
	"math"
	igroupoutlet "app/interface/group_outlet"
	"app/models"
	util "app/pkg/utils"
	"strings"
	"time"
	
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	
	uuid "github.com/satori/go.uuid"
)
	
type useGroupOutlet struct {
	repoGroupOutlet    igroupoutlet.Repository
	contextTimeOut time.Duration
}
	
func NewUseGroupOutlet(a igroupoutlet.Repository, timeout time.Duration) igroupoutlet.Usecase {
	return &useGroupOutlet{repoGroupOutlet: a, contextTimeOut: timeout}
}
	
func (u *useGroupOutlet) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.GroupOutlet, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	
	result, err = u.repoGroupOutlet.GetDataBy(ctx, ID)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useGroupOutlet) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	
	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}
	
	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoGroupOutlet.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}
	
	result.Total, err = u.repoGroupOutlet.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}
	
	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page
	
	return result, nil
}

func (u *useGroupOutlet) Create(ctx context.Context, Claims util.Claims, data *models.AddGroupOutlet) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mGroupOutlet = models.GroupOutlet{}
	)
	
	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mGroupOutlet.AddGroupOutlet)
	if err != nil {
		return err
	}

	mGroupOutlet.CreatedBy = uuid.FromStringOrNil(Claims.Id)
	mGroupOutlet.UpdatedBy = uuid.FromStringOrNil(Claims.Id)
	
	err = u.repoGroupOutlet.Create(ctx, &mGroupOutlet)
	if err != nil {
		return err
	}
	return nil
	
}

func (u *useGroupOutlet) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddGroupOutlet) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	
	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoGroupOutlet.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useGroupOutlet) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	
	err = u.repoGroupOutlet.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
	
		
	