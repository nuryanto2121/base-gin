package usetransaction

import (
	itransaction "app/interface/transaction"
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

type useTransaction struct {
	repoTransaction itransaction.Repository
	contextTimeOut  time.Duration
}

func NewUseTransaction(a itransaction.Repository, timeout time.Duration) itransaction.Usecase {
	return &useTransaction{repoTransaction: a, contextTimeOut: timeout}
}

func (u *useTransaction) GetDataBy(ctx context.Context, Claims util.Claims, ID uuid.UUID) (result *models.Transaction, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoTransaction.GetDataBy(ctx, "id", ID.String())
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *useTransaction) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoTransaction.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoTransaction.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useTransaction) Create(ctx context.Context, Claims util.Claims, data *models.TransactionForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mTransaction = models.Transaction{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mTransaction.AddTransaction)
	if err != nil {
		return err
	}

	mTransaction.CreatedBy = uuid.FromStringOrNil(Claims.Id)
	mTransaction.UpdatedBy = uuid.FromStringOrNil(Claims.Id)

	err = u.repoTransaction.Create(ctx, &mTransaction)
	if err != nil {
		return err
	}
	return nil

}

func (u *useTransaction) Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.TransactionForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoTransaction.Update(ctx, ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useTransaction) Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoTransaction.Delete(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}
