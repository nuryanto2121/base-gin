package useroleoutlet

import (
	iroleoutlet "app/interface/role_outlet"
	"app/models"
	util "app/pkg/util"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/jinzhu/copier"

	uuid "github.com/satori/go.uuid"
)

type useRoleOutlet struct {
	repoRoleOutlet iroleoutlet.Repository
	contextTimeOut time.Duration
}

func NewUseRoleOutlet(a iroleoutlet.Repository, timeout time.Duration) iroleoutlet.Usecase {
	return &useRoleOutlet{repoRoleOutlet: a, contextTimeOut: timeout}
}

func (u *useRoleOutlet) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.RoleOutlet, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoRoleOutlet.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useRoleOutlet) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoRoleOutlet.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoRoleOutlet.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useRoleOutlet) Create(ctx context.Context, Claims util.Claims, data *models.AddRoleOutlet) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mRoleOutlet = models.RoleOutlet{}
	)

	// mapping to struct model saRole
	err = copier.Copy(&mRoleOutlet.AddRoleOutlet, data)
	if err != nil {
		return err
	}

	mRoleOutlet.CreatedBy = uuid.FromStringOrNil(Claims.UserID)
	mRoleOutlet.UpdatedBy = uuid.FromStringOrNil(Claims.UserID)

	err = u.repoRoleOutlet.Create(ctx, &mRoleOutlet)
	if err != nil {
		return err
	}
	return nil

}

func (u *useRoleOutlet) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddRoleOutlet) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["updated_by"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoRoleOutlet.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useRoleOutlet) Delete(ctx context.Context, Claims util.Claims, UserId uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoRoleOutlet.Delete(ctx, "user_id", UserId.String())
	if err != nil {
		return err
	}
	return nil
}
