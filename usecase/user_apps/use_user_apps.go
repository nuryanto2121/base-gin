package useuserapps

import (
	itrx "app/interface/trx"
	iuserapps "app/interface/user_apps"
	"app/models"
	"app/pkg/logging"
	"app/pkg/util"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/structs"

	uuid "github.com/satori/go.uuid"
)

type useUserApps struct {
	repoUserApps   iuserapps.Repository
	repoTrx        itrx.Repository
	contextTimeOut time.Duration
}

func NewUseUserApps(a iuserapps.Repository, trx itrx.Repository, timeout time.Duration) iuserapps.Usecase {
	return &useUserApps{repoUserApps: a, repoTrx: trx, contextTimeOut: timeout}
}

func (u *useUserApps) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.UserApps, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoUserApps.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useUserApps) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
	}

	result.Data, err = u.repoUserApps.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoUserApps.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useUserApps) Create(ctx context.Context, Claims util.Claims, data *models.UserApps) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	// var (
	// 	mUserApps = models.UserApps{}
	// )

	// // mapping to struct model saRole
	// err = mapstructure.Decode(data, &mUserApps.AddUserApps)
	// if err != nil {
	// 	return err
	// }

	data.CreatedBy = uuid.FromStringOrNil(Claims.UserID)
	data.UpdatedBy = uuid.FromStringOrNil(Claims.UserID)

	err = u.repoUserApps.Create(ctx, data)
	if err != nil {
		return err
	}
	return nil

}

func (u *useUserApps) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.AddUserApps) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["created_by"] = Claims.UserID
	// fmt.Println(myMap)
	err = u.repoUserApps.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useUserApps) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUserApps.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}

// CreateChild implements iuserapps.Usecase
func (u *useUserApps) UpsertChild(ctx context.Context, Claims util.Claims, data models.ChildForm) (models.ChildForm, error) {
	ctxb, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		// wg     sync.WaitGroup
		logger = logging.Logger{}
		// ctxb = context.Background()
	)
	userID := uuid.FromStringOrNil(Claims.UserID)
	errTx := u.repoTrx.Run(ctxb, func(trxCtx context.Context) error {
		for _, val := range data.Childs {

			if val.ChildrenId == "" {
				//create
				Child := models.UserApps{
					AddUserApps: models.AddUserApps{
						Name:     val.Name,
						IsParent: false,
						ParentId: userID,
						DOB:      val.DOB,
					},
				}
				err := u.Create(ctxb, Claims, &Child)
				if err != nil {
					logger.Error("failed create child ", err)
					return err
				}
				val.ChildrenId = Child.Id.String()
			} else {
				//update
				ID := uuid.FromStringOrNil(val.ChildrenId)
				err := u.Update(ctxb, Claims, ID, &models.AddUserApps{
					Name:     val.Name,
					IsParent: false,
					ParentId: userID,
					DOB:      val.DOB,
				})
				if err != nil {
					logger.Error("failed update child ", err)
					return err
				}

			}

		}

		return nil
	})

	return data, errTx
}
