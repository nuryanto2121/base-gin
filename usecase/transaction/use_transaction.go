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
	"errors"
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
	// repuUserApps iuserapps.Repository
	contextTimeOut time.Duration
}

func NewUseTransaction(
	a itransaction.Repository,
	a1 itransactiondetail.Repository,
	b ioutlets.Repository,
	c iskumanagement.Repository,
	d iuserapps.Repository,
	e itrx.Repository, timeout time.Duration,
) itransaction.Usecase {
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

func (u *useTransaction) GetDataBy(ctx context.Context, Claims util.Claims, transactionId string) (*models.TransactionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var now = time.Now()

	//get data parent
	parent, err := u.repoCustomer.GetDataBy(ctx, "id", Claims.UserID)
	if err != nil {
		return nil, err
	}

	trxHeader, err := u.repoTransaction.GetDataBy(ctx, "transaction_code", transactionId)
	if err != nil {
		return nil, err
	}
	//get outlet
	outlet, err := u.repoOutlet.GetDataBy(ctx, "id", trxHeader.OutletId.String())
	if err != nil {
		return nil, err
	}

	//get detail
	trxDetail, err := u.repoTransDetail.GetList(ctx, models.ParamList{
		Page:       1,
		PerPage:    100,
		InitSearch: fmt.Sprintf("transaction_id = '%s'", trxHeader.Id),
	}) //GetDataBy(ctx, "transaction_id", trxHeader.Id.String())
	if err != nil {
		return nil, err
	}

	details := make([]*models.TransactionDetailResponse, 0, len(trxDetail))

	for _, val := range trxDetail {
		dt := &models.TransactionDetailResponse{}

		//get product
		product, err := u.repoSkuManagement.GetDataBy(ctx, "id", val.ProductId.String())
		if err != nil {
			return nil, err
		}

		dt.Description = fmt.Sprintf("%d x %s", val.ProductQty, product.SkuName)
		if val.IsChildren {
			child, err := u.repoCustomer.GetDataBy(ctx, "id", val.CustomerId.String())
			if err != nil {
				return nil, err
			}
			dt.CustomerName = child.Name
			dt.Description = fmt.Sprintf("%d x Durasi %d jam", val.ProductQty, product.Duration)
			//gen ticket
			if Claims.Role == "ticket" {

				child, err := u.repoCustomer.GetDataBy(ctx, "id", val.CustomerId.String())
				if err != nil {
					return nil, err
				}

				endTime := now.Add(time.Hour * time.Duration(val.Duration))
				qr := map[string]interface{}{
					"child_name":  child.Name,
					"parent_name": parent.Name,
					"phone_no":    parent.PhoneNo,
					"end_time":    endTime,
					"ticket_no":   val.TicketNo,
				}

				dt.QR = qr
			}
		}
		dt.Amount = val.Amount
		dt.Duration = val.Duration
		dt.ProductQty = val.ProductQty

		details = append(details, dt)
	}

	result := &models.TransactionResponse{
		TransactionId:         transactionId,
		TransactionDate:       trxHeader.TransactionDate,
		OutletName:            outlet.OutletName,
		OutletCity:            outlet.OutletCity,
		TotalTicket:           trxHeader.TotalTicket,
		TotalAmount:           trxHeader.TotalAmount,
		StatusTransaction:     trxHeader.StatusTransaction,
		StatusTransactionDesc: trxHeader.StatusTransaction.String(),
		StatusPayment:         trxHeader.StatusPayment,
		StatusPaymentDesc:     trxHeader.StatusPayment.String(),
		Details:               details,
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
	//get Outlet
	outlet, err := u.repoOutlet.GetDataBy(ctx, "id", req.OutletId.String())
	if err != nil {
		return nil, errors.New("outlets not found")
	}

	trxPrefix := fmt.Sprintf("BOK-%s", strings.ToUpper(outlet.OutletName[0:3]))
	t := &models.TmpCode{Prefix: trxPrefix}
	tsCode = util.GenCode(t)

	mTransaction := &models.Transaction{
		AddTransaction: models.AddTransaction{
			TransactionDate:   req.TransactionDate,
			OutletId:          req.OutletId,
			StatusPayment:     models.STATUS_WAITINGPAYMENT,
			StatusTransaction: models.STATUS_ORDER,
			TransactionCode:   tsCode,
			TotalAmount:       0,
			CustomerId:        uuid.FromStringOrNil(Claims.UserID),
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
			isChild := false
			customerName := ""
			Product, err := u.repoSkuManagement.GetDataBy(trxCtx, "id", val.ProductId.String())
			if err != nil {
				return err
			}

			Customer, err := u.repoCustomer.GetDataBy(trxCtx, "id", val.ChildrenId.String())
			if err != nil && err != models.ErrNotFound {
				return err
			}

			if Customer != nil {
				customerName = Customer.Name
			}
			//nnti diganti dengan outlet dan hari
			totalAmount += val.Amount //Product.PriceWeekday

			var desc = fmt.Sprintf("%d x %s", val.ProductQty, Product.SkuName)
			var TicketNo = ""
			if Product.IsBracelet {
				// t.Prefix = fmt.Sprintf("%s",)
				t.Prefix = fmt.Sprintf("TRC-%s", strings.ToUpper(outlet.OutletName[0:3]))
				TicketNo = util.GenCode(t)
				jmlTicket++
				isChild = true
				desc = fmt.Sprintf("%d x Durasi %d jam", val.ProductQty, Product.Duration)
			}

			//for response
			dtl = append(dtl, &models.TransactionDetailResponse{
				CustomerName: customerName,
				ProductQty:   val.ProductQty,
				Duration:     Product.Duration,
				Amount:       val.Amount,
				Description:  desc,
			})

			//insert to detail
			trxDetail := &models.TransactionDetail{
				AddTransactionDetail: models.AddTransactionDetail{
					TransactionId: mTransaction.Id,
					TicketNo:      TicketNo,
					CustomerId:    val.ChildrenId,
					IsChildren:    isChild,
					ProductId:     val.ProductId,
					ProductQty:    val.ProductQty,
					Amount:        val.Amount,
					Price:         val.Price,
					Duration:      Product.Duration,
				},
			}

			err = u.repoTransDetail.Create(trxCtx, trxDetail)
			if err != nil {
				return err
			}
		}

		// update header
		dtUpdate := map[string]interface{}{
			"total_amount": totalAmount,
			"total_ticket": jmlTicket,
		}
		return u.repoTransaction.Update(trxCtx, mTransaction.Id, dtUpdate)

	})
	if errTx != nil {
		return result, errTx
	}
	result.TotalTicket = jmlTicket
	result.TransactionDate = req.TransactionDate
	result.TotalAmount = totalAmount
	result.Details = dtl
	result.TransactionId = tsCode
	result.OutletCity = outlet.OutletCity
	result.OutletName = outlet.OutletName

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

// Payment implements itransaction.Usecase
func (u *useTransaction) Payment(ctx context.Context, Claims util.Claims, req *models.TransactionPaymentForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		now    = time.Now()
		userId = uuid.FromStringOrNil(Claims.UserID)
	)

	transaction, err := u.repoTransaction.GetDataBy(ctx, "transaction_code", req.TransactionId)
	if err != nil {
		return err
	}
	transaction.PaymentCode = req.PaymentCode
	transaction.Description = req.Description

	transaction.StatusPayment = models.STATUS_PAYMENTSUCCESS
	transaction.UpdatedAt = now
	transaction.UpdatedBy = userId

	err = u.repoTransaction.Update(ctx, transaction.Id, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (u *useTransaction) genTicket() string {
	var result = ""
	// 	Child Name:
	// Parents Name :
	// Phone Number :
	// End Time :
	// No tiket :
	return result
}
