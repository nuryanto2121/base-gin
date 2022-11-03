package usetransaction

import (
	iinventory "app/interface/inventory"
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
	"sync"
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
	useInventory    iinventory.Usecase
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
	g iinventory.Usecase,
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
		useInventory:      g,
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
	//validasi outlet active with trx
	if Claims.OutletId != trxHeader.OutletId.String() {
		return nil, models.ErrNoMatchOutlet
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
		SortField:  "ticket_no asc",
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
		if val.IsChildren || product.IsFree {
			child, err := u.repoCustomer.GetDataBy(ctx, "id", val.CustomerId.String())
			if err != nil {
				return nil, err
			}

			dt.CustomerName = child.Name
			dt.Description = fmt.Sprintf("%d x Durasi %d jam", val.ProductQty, product.Duration)
			if product.IsFree {
				dt.CustomerName = product.SkuName
				dt.Description = fmt.Sprintf("1 x %s", product.SkuName)
			}
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
				qr := map[string]interface{}{}
				if product.IsFree {
					qr = map[string]interface{}{
						"name":        parent.Name,
						"parent_name": product.SkuName,
						"phone_no":    parent.PhoneNo,
						"end_time":    endTime,
						"ticket_no":   val.TicketNo,
					}
				} else {
					qr = map[string]interface{}{
						"name":        child.Name,
						"parent_name": parent.Name,
						"phone_no":    parent.PhoneNo,
						"end_time":    endTime,
						"ticket_no":   val.TicketNo,
					}
				}

				dt.QR = qr
			}
		}
		dt.Amount = val.Amount
		dt.Duration = val.Duration
		dt.ProductId = val.ProductId
		dt.ProductQty = val.ProductQty
		dt.TicketNo = val.TicketNo
		if val.IsOvertime {
			dt.IsOvertime = val.IsOvertime
			dt.IsOvertimePaid = val.IsOvertimePaid
		}

		details = append(details, dt)
	}

	result := &models.TransactionResponse{
		ID:                    trxHeader.Id,
		TransactionCode:       trxHeader.TransactionCode,
		TransactionDate:       trxHeader.TransactionDate,
		OutletID:              outlet.Id,
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
		// queryparam.InitSearch += " and date(check_in) <> '0001-01-01'"
	}
	dataList, err := u.repoTransaction.GetList(ctx, queryparam)
	if err != nil {
		return result, err
	}

	for _, val := range dataList {

		val.StatusPaymentDesc = val.StatusPayment.String()
		val.StatusTransactionDesc = val.StatusTransaction.String()
		val.StatusTransactionDtlDesc = val.StatusTransactionDtl.String()
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
	listTicket, err := u.repoTransaction.GetListTicketUser(ctx, queryparam)
	if err != nil {
		return result, err
	}
	var (
		wg    sync.WaitGroup
		sLock = sync.RWMutex{}
		errc  = make(chan error)
	)

	for _, ticketDtl := range listTicket {
		details := []*models.TransactionDetailResponse{}

		wg.Add(1)

		go func(wgd *sync.WaitGroup, ticket *models.TransactionResponse, dtLock *sync.RWMutex, errch chan error) {
			defer wgd.Done()
			//get detail
			trxDetail, err := u.repoTransDetail.GetList(ctx, models.ParamList{
				Page:       1,
				PerPage:    100,
				InitSearch: fmt.Sprintf("transaction_id = '%s'", ticket.ID),
				SortField:  "ticket_no asc",
			}) //GetDataBy(ctx, "transaction_id", trxHeader.Id.String())
			if err != nil {
				errch <- err
				return
			}

			for _, val := range trxDetail {
				dt := &models.TransactionDetailResponse{}

				//get product
				product, err := u.repoSkuManagement.GetDataBy(ctx, "id", val.ProductId.String())
				if err != nil {
					errch <- err
				}

				dtLock.Lock()
				dt.Description = fmt.Sprintf("%d x %s", val.ProductQty, product.SkuName)
				if val.IsChildren || product.IsFree {
					child, err := u.repoCustomer.GetDataBy(ctx, "id", val.CustomerId.String())
					if err != nil {
						return
					}

					dt.CustomerName = child.Name
					dt.Description = fmt.Sprintf("%d x Durasi %d jam", val.ProductQty, product.Duration)
					if product.IsFree {
						dt.CustomerName = product.SkuName
						dt.Description = fmt.Sprintf("1 x %s", product.SkuName)
					}

				}
				dt.Amount = val.Amount
				dt.Duration = val.Duration
				dt.ProductId = val.ProductId
				dt.ProductQty = val.ProductQty
				dt.TicketNo = val.TicketNo
				if val.IsOvertime {
					dt.IsOvertime = val.IsOvertime
					dt.IsOvertimePaid = val.IsOvertimePaid
				}

				details = append(details, dt)
				dtLock.Unlock()
			}

			ticket.Details = details

		}(&wg, ticketDtl, &sLock, errc)

	}

	wg.Wait()

	result.Data = listTicket

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
		var ProductAdultFree = &models.OutletList{}
		//create trx header
		err := u.repoTransaction.Create(trxCtx, &mTransaction)
		if err != nil {
			return err
		}
		//set id transaction
		result.ID = mTransaction.Id
		for _, val := range req.Details {
			//get sku adult free
			if val.IsAdultFree {
				productAdultFree, err := u.repoOutlet.GetList(trxCtx, models.ParamList{
					Page:    1,
					PerPage: 1,
					InitSearch: fmt.Sprintf(`
					o.id='%s' and 
					is_free=true and
					is_bracelet=true
					`, req.OutletId),
				})
				if err != nil {
					return err
				}
				if len(productAdultFree) > 0 {
					ProductAdultFree = productAdultFree[0]
					//check stock adult free then update inventory
					err := u.useInventory.PatchStock(trxCtx, Claims, models.InvPatchStockRequest{
						OutletId:  req.OutletId,
						ProductId: ProductAdultFree.ProductId,
						Qty:       -val.ProductQty,
					})
					if err != nil {
						return err
					}
				}

			}
			//check stock then update inventory
			err := u.useInventory.PatchStock(trxCtx, Claims, models.InvPatchStockRequest{
				OutletId:  req.OutletId,
				ProductId: val.ProductId,
				Qty:       -val.ProductQty,
			})
			if err != nil {
				return err
			}

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
			totalAmount += val.Amount //Product.Price

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
				ProductId:    val.ProductId,
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

			if val.IsAdultFree {
				// ProductAdultFree
				desc := fmt.Sprintf("1 x %s", ProductAdultFree.SkuName)
				dtl = append(dtl, &models.TransactionDetailResponse{
					CustomerName: ProductAdultFree.SkuName,
					ProductQty:   1,
					Duration:     ProductAdultFree.Duration,
					ProductId:    ProductAdultFree.ProductId,
					Amount:       0,
					Description:  desc,
				})

				//insert to detail
				trxDetail := &models.TransactionDetail{
					AddTransactionDetail: models.AddTransactionDetail{
						TransactionId: mTransaction.Id,
						TicketNo:      TicketNo + "/FREE",
						CustomerId:    val.ChildrenId,
						IsChildren:    false,
						ProductId:     ProductAdultFree.ProductId,
						ProductQty:    1,
						Amount:        0,
						Price:         0,
						Duration:      ProductAdultFree.Duration,
					},
				}

				err = u.repoTransDetail.Create(trxCtx, trxDetail)
				if err != nil {
					return err
				}
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
		now                = time.Now()
		logger             = logging.Logger{}
		userId             = uuid.FromStringOrNil(Claims.UserID)
		isAllCheckOut bool = true
	)

	//get data transaction
	transaction, err := u.repoTransaction.GetDataBy(ctx, "id", req.TransactionId)
	if err != nil {
		return nil, err
	}

	if transaction.Id == uuid.Nil {
		return nil, models.ErrTransactionNotFound
	}

	//validasi outlet active with trx
	if Claims.OutletId != transaction.OutletId.String() {
		return nil, models.ErrNoMatchOutlet
	}

	transaction.PaymentCode = req.PaymentCode
	transaction.Description = req.Description

	// if req.PaymentCode == models.PAYMENT_CASH {
	transaction.StatusPayment = models.STATUS_PAYMENTSUCCESS
	// } else {
	// 	transaction.StatusPayment = models.STATUS_WAITINGPAYMENT
	// }
	errTx := u.repoTrx.Run(ctx, func(trxCtx context.Context) error {
		//pembayaran overtime
		if transaction.StatusTransaction == models.STATUS_OVERTIME {
			trxDtl, err := u.repoTransaction.GetList(trxCtx, models.ParamList{
				Page:       1,
				PerPage:    10,
				InitSearch: fmt.Sprintf("t.id='%s' and td.is_children = true", transaction.Id),
			})
			if err != nil {
				logger.Error("error get transaction top one ", err)
				return err
			}

			isAllCheckOut = true

			for _, val := range trxDtl {
				if val.CheckOut.IsZero() {
					// if val.TicketNo != req.TicketNo {
					isAllCheckOut = false
					// }
				}
			}

			if isAllCheckOut {
				transaction.StatusTransaction = models.STATUS_FINISH
			} else {
				transaction.StatusTransaction = models.STATUS_ACTIVE
			}
			//update trx detail is_ov_paid true
			for _, val := range req.TicketOvertime {

				// trxDtl.AddTransactionDetail.IsOvertimePaid = true
				dtlFree := map[string]interface{}{
					"is_overtime_paid": true,
				}

				err = u.repoTransDetail.UpdateBy(trxCtx, fmt.Sprintf("ticket_no='%s' and is_overtime=true", val), dtlFree)
				if err != nil {
					logger.Error("error update ticket overtime ", err)
					return err
				}

			}

		} else {
			transaction.StatusTransaction = models.STATUS_ORDER
		}

		transaction.UpdatedAt = now
		transaction.UpdatedBy = userId

		err = u.repoTransaction.Update(trxCtx, transaction.Id, transaction)
		if err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return nil, errTx
	}

	return result, nil
}

// CheckIn implements itransaction.Usecase
func (u *useTransaction) CheckIn(ctx context.Context, Claims util.Claims, req *models.CheckInCheckOutForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var logger = logging.Logger{}

	transactionDetail, err := u.repoTransDetail.GetDataBy(ctx, "ticket_no", req.TicketNo)
	if err != nil {
		logger.Error("error get trx detail by ticket no ", req.TicketNo)
		return err
	}

	//getHeader
	transaction, err := u.repoTransaction.GetDataBy(ctx, "id", transactionDetail.AddTransactionDetail.TransactionId.String())
	if err != nil {
		logger.Error("error get trx by id ", transactionDetail.AddTransactionDetail.TransactionId)
		return err
	}
	if transaction.Id == uuid.Nil {
		return models.ErrTransactionNotFound
	}

	//validasi outlet active with trx
	if Claims.OutletId != transaction.OutletId.String() {
		return models.ErrNoMatchOutlet
	}

	if transaction.StatusPayment != models.STATUS_PAYMENTSUCCESS {
		return models.ErrPaymentNeeded
	}

	if transaction.StatusTransaction != models.STATUS_ORDER {
		return models.ErrNoStatusOrder
	}

	errTx := u.repoTrx.Run(ctx, func(trxCtx context.Context) error {

		trxDtl, err := u.repoTransaction.GetList(trxCtx, models.ParamList{
			Page:       1,
			PerPage:    10,
			InitSearch: fmt.Sprintf("t.id='%s' and td.is_children = true", transaction.Id),
		})
		if err != nil {
			logger.Error("error get list trx ", transaction.Id)
			return err
		}

		isCheckin := true

		for _, val := range trxDtl {
			if val.CheckIn.IsZero() {
				if val.TicketNo != req.TicketNo {
					isCheckin = false
				}
			}
		}

		transactionDetail.CheckIn = req.CheckIn
		transactionDetail.StatusTransactionDtl = models.STATUS_CHECKIN
		err = u.repoTransDetail.Update(trxCtx, transactionDetail.Id, transactionDetail)
		if err != nil {
			logger.Error("error update trx detail ", err)
			return err
		}

		// go func() {
		dtlFree := map[string]interface{}{
			"check_in":               req.CheckIn,
			"status_transaction_dtl": models.STATUS_CHECKIN,
		}
		err = u.repoTransDetail.UpdateBy(trxCtx, fmt.Sprintf("ticket_no='%s/FREE'", transactionDetail.TicketNo), dtlFree)
		if err != nil {
			logger.Error("error update trx detail free", err)
		}
		// }()

		if isCheckin {
			if transaction.StatusTransaction != models.STATUS_ACTIVE {
				transaction.StatusTransaction = models.STATUS_ACTIVE
				err = u.repoTransaction.Update(trxCtx, transaction.Id, transaction)
				if err != nil {
					logger.Error("error update trx ", err)
					return err
				}
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

	var (
		logger         = logging.Logger{}
		isCheckOut     bool
		isChildStillOV bool = false
	)

	transactionDetail, err := u.repoTransDetail.GetDataBy(ctx, "ticket_no", req.TicketNo)
	if err != nil {
		logger.Error("error get transaction detail", err)
		return err
	}

	//getHeader
	transaction, err := u.repoTransaction.GetDataBy(ctx, "id", transactionDetail.AddTransactionDetail.TransactionId.String())
	if err != nil {
		logger.Error("error get transaction ", err)
		return err
	}

	if transaction.Id == uuid.Nil {
		return models.ErrTransactionNotFound
	}

	//validasi outlet active with trx
	if Claims.OutletId != transaction.OutletId.String() {
		return models.ErrNoMatchOutlet
	}

	// if transaction.StatusPayment != models.STATUS_PAYMENTSUCCESS {
	// 	return models.ErrBadParamInput
	// }

	if !transactionDetail.CheckOut.IsZero() {
		return models.ErrNoStatusCheckIn
	}

	errTx := u.repoTrx.Run(ctx, func(trxCtx context.Context) error {
		trxDtl, err := u.repoTransaction.GetList(trxCtx, models.ParamList{
			Page:       1,
			PerPage:    10,
			InitSearch: fmt.Sprintf("t.id='%s' and td.is_children = true", transaction.Id),
		})
		if err != nil {
			logger.Error("error get transaction list ", transaction.Id, err)
			return err
		}

		isCheckOut = true

		for _, val := range trxDtl {
			if val.CheckOut.IsZero() {
				if val.TicketNo != req.TicketNo {
					isCheckOut = false
				}

			}
			if val.IsOvertime && !val.IsOvertimePaid {
				isChildStillOV = true
			}
		}

		transactionDetail.CheckOut = req.CheckOut //.In(util.GetLocation())
		transactionDetail.StatusTransactionDtl = models.STATUS_CHECKOUT
		err = u.repoTransDetail.Update(trxCtx, transactionDetail.Id, transactionDetail)
		if err != nil {
			logger.Error("error update transaction detail", err)
			return err
		}

		// go func() {
		dtlFree := map[string]interface{}{
			"check_out":              req.CheckOut,
			"status_transaction_dtl": models.STATUS_CHECKOUT,
		}
		err = u.repoTransDetail.UpdateBy(trxCtx, fmt.Sprintf("ticket_no='%s/FREE'", transactionDetail.TicketNo), dtlFree)
		if err != nil {
			logger.Error("error update trx detail free", err)
		}
		// }()

		if isCheckOut {
			if transaction.StatusTransaction != models.STATUS_FINISH {
				transaction.StatusTransaction = models.STATUS_FINISH
				if isChildStillOV {
					transaction.StatusTransaction = models.STATUS_OVERTIME
				}

				u.repoTransaction.Update(trxCtx, transaction.Id, transaction)
				if err != nil {
					logger.Error("error update transaction ", err)
					return err
				}
			}
		}

		return nil
	})

	if errTx != nil {
		return errTx
	}

	//validasi overtime
	price, overtimeAmt, hour, errOv := u.Overtime(ctx, transactionDetail, req.CheckOut, transaction.OutletId.String())
	if errOv != nil && overtimeAmt > 0 {
		ctxOv := context.Background()
		go u.repoTrx.Run(ctxOv, func(trxCtx context.Context) error {
			transactionDetail.Id = uuid.Nil
			transactionDetail.IsOvertime = true
			transactionDetail.Amount = overtimeAmt
			transactionDetail.Price = price
			transactionDetail.Duration = hour
			transactionDetail.StatusTransactionDtl = models.STATUS_OVERTIME
			transactionDetail.IsOvertimePaid = false
			err = u.repoTransDetail.Create(trxCtx, transactionDetail)
			if err != nil {
				logger.Error("error create transaction detail overtime", err)
				return err
			}

			transaction.TotalAmount += overtimeAmt
			// if isCheckOut { //if can partial u can uncomment this validation
			//if all child has checkout and have overtime then status payment waiting payment
			transaction.StatusPayment = models.STATUS_WAITINGPAYMENT
			// }
			transaction.StatusTransaction = models.STATUS_OVERTIME
			err = u.repoTransaction.Update(trxCtx, transaction.Id, transaction)
			if err != nil {
				logger.Error("error update transaction amount overtime", err)
				return err
			}
			return nil
		})
		//

		return errOv
	}

	if isChildStillOV {
		return models.ErrOvertime
	}

	return nil
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

func (u *useTransaction) Overtime(ctx context.Context, trxDetail *models.TransactionDetail, checkOut time.Time, outletId string) (price float64, amount float64, hour int64, err error) {
	var (
		// now = util.GetTimeNow()
		logger = logging.Logger{}
	)
	//getOutlet for check overtime
	outlet, err := u.repoOutlet.GetDataBy(ctx, "id", outletId)
	if err != nil {
		logger.Error("error get outlet ", err)
		return price, amount, hour, err
	}

	//check overtime
	diff := checkOut.Sub(trxDetail.CheckIn).Minutes() //now.Sub(checkOut).Minutes() //req.CheckOut.Sub(now).Minutes()
	timeOut := math.Round(diff)
	hour = int64(math.Ceil(float64(timeOut) / float64(60)))

	hour -= trxDetail.Duration
	if outlet.ToleransiTime > 0 && outlet.ToleransiTime < int64(timeOut-60) {

		overtimeAmount := float64(hour) * outlet.OvertimeAmount
		return outlet.OvertimeAmount, overtimeAmount, hour, models.ErrOvertime
	}
	return price, amount, hour, nil
}
