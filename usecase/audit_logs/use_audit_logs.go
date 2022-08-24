package useauditlogs

import (
	iauditlogs "app/interface/audit_logs"
	"app/models"
	"app/pkg/util"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"

	uuid "github.com/satori/go.uuid"
)

type useAuditLogs struct {
	repoAuditLogs  iauditlogs.Repository
	contextTimeOut time.Duration
}

func NewUseAuditLogs(a iauditlogs.Repository, timeout time.Duration) iauditlogs.Usecase {
	return &useAuditLogs{repoAuditLogs: a, contextTimeOut: timeout}
}

func (u *useAuditLogs) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.AuditLogs, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoAuditLogs.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useAuditLogs) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoAuditLogs.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoAuditLogs.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useAuditLogs) Create(ctx context.Context, Claims util.Claims, data *models.AddAuditLogs) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mAuditLogs = models.AuditLogs{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mAuditLogs.AddAuditLogs)
	if err != nil {
		return err
	}

	mAuditLogs.CreatedBy = uuid.FromStringOrNil(Claims.Id)
	mAuditLogs.UpdatedBy = uuid.FromStringOrNil(Claims.Id)

	err = u.repoAuditLogs.Create(ctx, &mAuditLogs)
	if err != nil {
		return err
	}
	return nil

}

func (u *useAuditLogs) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddAuditLogs) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoAuditLogs.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useAuditLogs) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoAuditLogs.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
