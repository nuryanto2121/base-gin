package usetransaction

import (
	ioutlets "app/interface/outlets"
	iskumanagement "app/interface/sku_management"
	itransaction "app/interface/transaction"
	itransactiondetail "app/interface/transaction_detail"
	itrx "app/interface/trx"
	iuserapps "app/interface/user_apps"
	"app/models"
	"app/pkg/util"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/structs"

	uuid "github.com/satori/go.uuid"
)

type useTransaction struct {
	repoTransaction   itransaction.Repository
	repoTransDetail   itransactiondetail.Repository
	repoOutlet        ioutlets.Repository
	repoSkuManagement iskumanagement.Repository
	repoCustomer      iuserapps.Repository
	repoTrx           itrx.Repository
	contextTimeOut    time.Duration
}

func NewUseTransaction(a itransaction.Repository, a1 itransactiondetail.Repository, b ioutlets.Repository, c iskumanagement.Repository, d iuserapps.Repository, e itrx.Repository, timeout time.Duration) itransaction.Usecase {
	return &useTransaction{
		repoTransaction:   a,
		repoTransDetail:   a1,
		repoOutlet:        b,
		repoSkuManagement: c,
		repoCustomer:      d,
		repoTrx:           e,
		contextTimeOut:    timeout,
	}
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

func (u *useTransaction) Create(ctx context.Context, Claims util.Claims, req *models.TransactionForm) (*models.TransactionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		// mTransaction = models.Transaction{}
		tsCode      string  = ""
		jmlTicket   int64   = 0
		result              = &models.TransactionResponse{}
		dtl                 = []*models.TransactionDetailResponse{}
		totalAmount float64 = 0
	)
	t := &models.TmpCode{Prefix: "TRX"}
	tsCode = util.GenCode(t)

	mTransaction := &models.Transaction{
		AddTransaction: models.AddTransaction{
			TransactionDate:   req.TransactionDate,
			OutletId:          req.OutletId,
			StatusPayment:     models.STATUS_WAITINGPAYMENT,
			StatusTransaction: models.STATUS_ORDER,
			TransactionCode:   tsCode,
			TotalAmount:       0,
		},
		Model: models.Model{
			CreatedBy: uuid.FromStringOrNil(Claims.UserID),
			UpdatedBy: uuid.FromStringOrNil(Claims.UserID),
		},
	}

	errTx := u.repoTrx.Run(ctx, func(trxCtx context.Context) error {
		err := u.repoTransaction.Create(trxCtx, mTransaction)
		if err != nil {
			return err
		}

		for _, val := range req.Details {
			Product, err := u.repoSkuManagement.GetDataBy(trxCtx, "id", val.ProductId.String())
			if err != nil {
				return err
			}

			Customer, err := u.repoCustomer.GetDataBy(trxCtx, "id", val.CustomerId.String())
			if err != nil {
				return err
			}

			trxDetail := &models.TransactionDetail{
				AddTransactionDetail: models.AddTransactionDetail{
					TransactionId: mTransaction.Id,
					CustomerId:    val.CustomerId,
					IsParent:      false,
					ProductId:     val.ProductId,
					ProductQty:    val.ProductQty,
					Amount:        Product.PriceWeekday,
					Duration:      Product.Duration,
				},
			}

			err = u.repoTransDetail.Create(trxCtx, trxDetail)
			if err != nil {
				return err
			}
			//nnti diganti dengan outlet dan hari
			totalAmount += Product.PriceWeekday

			var desc = fmt.Sprintf("%d x %s", val.ProductQty, Product.SkuName)
			if Product.IsBracelet {
				jmlTicket++
				desc = fmt.Sprintf("%d x Durasi %d jam", val.ProductQty, Product.Duration)
			}

			dtl = append(dtl, &models.TransactionDetailResponse{
				CustomerName: Customer.Name,
				ProductQty:   val.ProductQty,
				Duration:     Product.Duration,
				Amount:       Product.PriceWeekday,
				Description:  desc,
			})
		}

		// update header
		dtUpdate := map[string]interface{}{
			"total_amount": totalAmount,
		}
		u.repoTransaction.Update(trxCtx, mTransaction.Id, dtUpdate)
		return nil
	})
	if errTx != nil {
		return result, errTx
	}
	result.TotalTicket = jmlTicket
	result.TransactionDate = req.TransactionDate
	result.TotalAmount = totalAmount
	result.Details = dtl

	return result, nil

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
