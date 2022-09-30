package usetransaction

import (
	ioutlets "app/interface/outlets"
	iskumanagement "app/interface/sku_management"
	itransaction "app/interface/transaction"
	itransactiondetail "app/interface/transaction_detail"
	itrx "app/interface/trx"
	iuserapps "app/interface/user_apps"
	"app/models"
	"app/pkg/logging"
	_midtransGateway "app/pkg/midtrans"
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
	midtransGateway *_midtransGateway.Gateway
	contextTimeOut  time.Duration
}

func NewUseTransaction(
	a itransaction.Repository,
	a1 itransactiondetail.Repository,
	b ioutlets.Repository,
	c iskumanagement.Repository,
	d iuserapps.Repository,
	e itrx.Repository,
	f *_midtransGateway.Gateway,
	timeout time.Duration,
) itransaction.Usecase {
	return &useTransaction{
		repoTransaction:   a,
		repoTransDetail:   a1,
		repoOutlet:        b,
		repoSkuManagement: c,
		repoCustomer:      d,
		repoTrx:           e,
		midtransGateway:   f,
		contextTimeOut:    timeout,
	}
}

func (u *useTransaction) GetDataBy(ctx context.Context, Claims util.Claims, transactionId string) (*models.TransactionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		now    = time.Now()
		parent = &models.UserApps{}
		err    error
	)

	sField := "id"
	if Claims.Id == "transactionCode" {
		sField = "transaction_code"
	}

	trxHeader, err := u.repoTransaction.GetDataBy(ctx, sField, transactionId)
	if err != nil {
		return nil, err
	}
	if trxHeader.Id == uuid.Nil {
		return nil, models.ErrTransactionNotFound
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

				//get parent
				if parent.Name == "" {
					parent, err = u.repoCustomer.GetDataBy(ctx, "id", child.ParentId.String())
					if err != nil {
						return nil, err
					}
				}

				child, err := u.repoCustomer.GetDataBy(ctx, "id", val.CustomerId.String())
				if err != nil {
					return nil, err
				}

				endTime := now.Add(time.Hour * time.Duration(val.Duration))
				qr := map[string]interface{}{
					"name":        child.Name,
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
		ID:                    trxHeader.Id,
		TransactionCode:       trxHeader.TransactionCode,
		TransactionDate:       trxHeader.TransactionDate,
		OutletName:            outlet.OutletName,
		OutletCity:            outlet.OutletCity,
		TotalTicket:           trxHeader.TotalTicket,
		TotalAmount:           trxHeader.TotalAmount,
		StatusTransaction:     trxHeader.StatusTransaction,
		StatusTransactionDesc: trxHeader.StatusTransaction.String(),
		StatusPayment:         trxHeader.StatusPayment,
		StatusPaymentDesc:     trxHeader.StatusPayment.String(),
		PaymentId:             trxHeader.PaymentId,
		PaymentToken:          trxHeader.PaymentToken,
		PaymentStatusDesc:     trxHeader.PaymentStatusDesc,
		PaymentType:           trxHeader.Description,
		PaymentCode:           trxHeader.PaymentCode,
		PaymentCodeDesc:       trxHeader.PaymentCode.String(),
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
		queryparam.InitSearch += " and date(check_in) <> '0001-01-01'"
	}
	dataList, err := u.repoTransaction.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	for _, val := range dataList {

		val.StatusPaymentDesc = val.StatusPayment.String()
		val.StatusTransactionDesc = val.StatusTransaction.String()
		// if val.StatusTransaction == models.STATUS_CHECKIN { //!val.CheckIn.IsZero() {
		// 	// fmt.Println(val.CheckIn.In(util.GetLocation()))
		// 	// fmt.Println("=======================")
		// 	// fmt.Println(val.CheckIn)
		// 	// fmt.Println("")
		// fmt.Println(util.GetTimeNow())
		// 	// fmt.Println("")
		// 	if val.CheckIn.Before(util.GetTimeNow()) {
		// 		fmt.Println(val.CheckIn)
		// 		fmt.Println("sebelum")
		// 		fmt.Println(util.GetTimeNow())
		// 		fmt.Println("-----------")
		// 	}

		// 	if val.CheckIn.After(util.GetTimeNow()) {
		// 		fmt.Println(val.CheckIn)
		// 		fmt.Println("Sesudah")
		// 		fmt.Println(util.GetTimeNow())
		// 		fmt.Println("=============")
		// 	}
		// upt := map[string]interface{}{
		// 	"check_in": util.GetTimeNow(),
		// }

		// err := u.repoTransDetail.UpdateBy(ctx, fmt.Sprintf("ticket_no = '%s'", val.TicketNo), upt)
		// if err != nil {
		// 	fmt.Printf("%v", err)
		// }
		// }

	}

	result.Data = dataList

	result.Total, err = u.repoTransaction.Count(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int64(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useTransaction) GetListTicketUser(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += fmt.Sprintf(" AND t.customer_id = '%s'", Claims.UserID)
	} else {
		queryparam.InitSearch = fmt.Sprintf("t.customer_id = '%s'", Claims.UserID)
	}
	result.Data, err = u.repoTransaction.GetListTicketUser(ctx, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoTransaction.CountUserList(ctx, queryparam)
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

	//check has have staus draf
	pWhere := fmt.Sprintf("customer_id = '%s' and status_transaction = %d", Claims.UserID, models.STATUS_DRAF)
	fctExist, err := u.repoTransaction.IsExist(ctx, pWhere)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}
	if fctExist {
		return nil, models.ErrStillHaveDraf
	}
	//get Outlet
	outlet, err := u.repoOutlet.GetDataBy(ctx, "id", req.OutletId.String())
	if err != nil {
		return nil, errors.New("outlets not found")
	}

	trxPrefix := fmt.Sprintf("BOK-%s", strings.ToUpper(outlet.OutletName[0:3]))
	t := &models.TmpCode{Prefix: trxPrefix}
	tsCode = util.GenCode(t)

	mTransaction := models.Transaction{
		AddTransaction: models.AddTransaction{
			TransactionDate:   req.TransactionDate,
			OutletId:          req.OutletId,
			StatusPayment:     models.STATUS_WAITINGPAYMENT,
			StatusTransaction: models.STATUS_DRAF,
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
		err := u.repoTransaction.Create(trxCtx, &mTransaction)
		if err != nil {
			return err
		}
		//set id transaction
		result.ID = mTransaction.Id
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
				Prefix := fmt.Sprintf("TRC-%s", strings.ToUpper(outlet.OutletName[0:3]))
				t := &models.TmpCode{Prefix: Prefix}
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
	// result.ID
	result.TotalTicket = jmlTicket
	result.TransactionDate = req.TransactionDate
	result.TotalAmount = totalAmount
	result.Details = dtl
	result.TransactionCode = tsCode
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

	errTx := u.repoTrx.Run(ctx, func(trxCtx context.Context) error {
		err = u.repoTransaction.Delete(trxCtx, ID)
		if err != nil {
			return err
		}

		err = u.repoTransDetail.Delete(trxCtx, ID)
		if err != nil {
			return err
		}
		return nil
	})

	return errTx
}

// Payment implements itransaction.Usecase
func (u *useTransaction) Payment(ctx context.Context, Claims util.Claims, req *models.TransactionPaymentForm) (result *models.MidtransResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		now = time.Now()
		// logger = logging.Logger{}
		userId = uuid.FromStringOrNil(Claims.UserID)
	)

	//get data transaction
	transaction, err := u.repoTransaction.GetDataBy(ctx, "id", req.TransactionId)
	if err != nil {
		return nil, err
	}

	if transaction.Id == uuid.Nil {
		return nil, models.ErrTransactionNotFound
	}

	transaction.PaymentCode = req.PaymentCode
	transaction.Description = req.Description

	// if req.PaymentCode == models.PAYMENT_CASH {
	transaction.StatusPayment = models.STATUS_PAYMENTSUCCESS
	// } else {
	// 	transaction.StatusPayment = models.STATUS_WAITINGPAYMENT
	// }

	transaction.StatusTransaction = models.STATUS_ORDER

	transaction.UpdatedAt = now
	transaction.UpdatedBy = userId

	err = u.repoTransaction.Update(ctx, transaction.Id, transaction)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CheckIn implements itransaction.Usecase
func (u *useTransaction) CheckIn(ctx context.Context, Claims util.Claims, req *models.CheckInCheckOutForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	transactionDetail, err := u.repoTransDetail.GetDataBy(ctx, "ticket_no", req.TicketNo)
	if err != nil {
		return err
	}

	//getHeader
	transaction, err := u.repoTransaction.GetDataBy(ctx, "id", transactionDetail.AddTransactionDetail.TransactionId.String())
	if err != nil {
		return err
	}
	if transaction.Id == uuid.Nil {
		return models.ErrTransactionNotFound
	}

	if transaction.StatusPayment != models.STATUS_PAYMENTSUCCESS {
		return models.ErrPaymentNeeded
	}

	if transaction.StatusTransaction != models.STATUS_ORDER {
		return models.ErrNoStatusOrder
	}

	errTx := u.repoTrx.Run(ctx, func(trxCtx context.Context) error {

		trxDtl, err := u.repoTransDetail.GetList(ctx, models.ParamList{
			Page:       1,
			PerPage:    10,
			InitSearch: fmt.Sprintf("parent_id='%s' and is_child = true", transaction.Id),
		})
		if err != nil {
			return err
		}

		isCheckin := true

		for _, val := range trxDtl {
			if val.CheckIn.IsZero() {
				isCheckin = false
			}
		}

		if isCheckin {
			transactionDetail.CheckIn = req.CheckIn.In(util.GetLocation())
			err = u.repoTransDetail.Update(trxCtx, transactionDetail.Id, transactionDetail)
			if err != nil {
				return err
			}
		}

		if transaction.StatusTransaction != models.STATUS_CHECKIN {
			transaction.StatusTransaction = models.STATUS_CHECKIN
			u.repoTransaction.Update(trxCtx, transaction.Id, transaction)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return errTx
}

// CheckOut implements itransaction.Usecase
func (u *useTransaction) CheckOut(ctx context.Context, Claims util.Claims, req *models.CheckInCheckOutForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	transactionDetail, err := u.repoTransDetail.GetDataBy(ctx, "ticket_no", req.TicketNo)
	if err != nil {
		return err
	}

	//getHeader
	transaction, err := u.repoTransaction.GetDataBy(ctx, "id", transactionDetail.AddTransactionDetail.TransactionId.String())
	if err != nil {
		return err
	}

	if transaction.Id == uuid.Nil {
		return models.ErrTransactionNotFound
	}

	if transaction.StatusPayment != models.STATUS_PAYMENTSUCCESS {
		return models.ErrBadParamInput
	}

	// if transaction.StatusTransaction != models.STATUS_CHECKIN {
	// 	return models.ErrNoStatusCheckIn
	// }
	if !transactionDetail.CheckOut.IsZero() {
		return models.ErrNoStatusCheckIn
	}

	errTx := u.repoTrx.Run(ctx, func(trxCtx context.Context) error {
		trxDtl, err := u.repoTransDetail.GetList(ctx, models.ParamList{
			Page:       1,
			PerPage:    10,
			InitSearch: fmt.Sprintf("parent_id='%s' and is_child = true", transaction.Id),
		})
		if err != nil {
			return err
		}

		isCheckOut := true

		for _, val := range trxDtl {
			if val.CheckOut.IsZero() {
				isCheckOut = false
			}
		}

		if isCheckOut {
			transactionDetail.CheckOut = req.CheckOut.In(util.GetLocation())
			err = u.repoTransDetail.Update(ctx, transactionDetail.Id, transactionDetail)
			if err != nil {
				return err
			}
		}

		if transaction.StatusTransaction != models.STATUS_CHECKOUT {
			transaction.StatusTransaction = models.STATUS_CHECKOUT
			u.repoTransaction.Update(trxCtx, transaction.Id, transaction)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return errTx
}

// GetById implements itransaction.Usecase
func (u *useTransaction) GetById(ctx context.Context, Claims util.Claims, transactionId string) (result *models.Transaction, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var logger = logging.Logger{}

	result, err = u.repoTransaction.GetDataBy(ctx, "id", transactionId)
	if err != nil {
		logger.Error("error can't find transaction by id ", transactionId)
		return nil, err
	}
	return result, nil
}

// UpdateHeader implements itransaction.Usecase
func (u *useTransaction) UpdateHeader(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.Transaction) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	//check data is exist
	isExist, err := u.repoTransaction.IsExist(ctx, fmt.Sprintf("id = '%s'", data.Id))
	if err != nil {
		return models.ErrNotFound
	}

	if !isExist {
		return models.ErrDataAlreadyExist
	}

	data.UpdatedBy = uuid.FromStringOrNil(Claims.UserID)
	err = u.repoTransaction.Update(ctx, ID, &data)
	if err != nil {
		return err
	}
	return nil
}
