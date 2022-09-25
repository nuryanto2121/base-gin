package usepayment

import (
	imidtrans "app/interface/midtrans"
	ipayment "app/interface/payment"
	ipaymentnotificationlogs "app/interface/payment_notification_logs"
	itransaction "app/interface/transaction"
	"app/models"
	"app/pkg/logging"
	"app/pkg/util"
	"context"
	"fmt"
	"time"

	_midtransGateway "app/pkg/midtrans"

	// "github.com/midtrans/midtrans-go"

	"github.com/midtrans/midtrans-go"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type usePayment struct {
	useTransaction      itransaction.Usecase
	repoTransaction     itransaction.Repository
	repoPaymentNotifLog ipaymentnotificationlogs.Repository
	midtransGateway     *_midtransGateway.Gateway
	// coreGateway         *_coreGateway.Gateway
	coreGateway    imidtrans.Repository
	contextTimeOut time.Duration
}

func NewUsePayment(a itransaction.Usecase, b *_midtransGateway.Gateway, c itransaction.Repository, d ipaymentnotificationlogs.Repository, e imidtrans.Repository, timeout time.Duration) ipayment.Usecase {
	return &usePayment{
		useTransaction:      a,
		midtransGateway:     b,
		repoTransaction:     c,
		repoPaymentNotifLog: d,
		coreGateway:         e,
		contextTimeOut:      timeout,
	}
}

// Payment implements ipayment.Usecase
func (u *usePayment) Payment(ctx context.Context, Claims util.Claims, req *models.TransactionPaymentForm) (result *models.MidtransResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		now    = time.Now()
		logger = logging.Logger{}
		userId = uuid.FromStringOrNil(Claims.UserID)
	)

	//get data transaction
	transaction, err := u.useTransaction.GetById(ctx, Claims, req.TransactionId)
	if err != nil {
		return nil, err
	}

	if transaction.Id == uuid.Nil {
		return nil, models.ErrTransactionNotFound
	}

	transaction.PaymentCode = req.PaymentCode
	transaction.Description = req.Description

	if req.PaymentCode == models.PAYMENT_CASH {
		transaction.StatusPayment = models.STATUS_PAYMENTSUCCESS
	} else {
		transaction.StatusPayment = models.STATUS_WAITINGPAYMENT
	}

	transaction.StatusTransaction = models.STATUS_ORDER

	if req.PaymentCode != models.PAYMENT_CASH && req.PaymentCode != models.PAYMENT_CASHIER {
		// generate request midtrans
		invBuilder, err := u.useTransaction.BuildMidtrans(ctx, transaction)
		if err != nil {
			return nil, err
		}

		reqSnap, err := invBuilder.Build()
		if err != nil {
			return nil, err
		}

		//hit payment midtrans
		res, err := u.midtransGateway.SnapV2Gateway.CreateTransaction(reqSnap)
		e, ok := err.(*midtrans.Error)

		if ok && e != nil {
			logger.Error("error hit to midranst ", e)
			return nil, e
		}
		err = mapstructure.Decode(res, &result)
		if err != nil {
			return nil, err
		}

		//update payment token header trx
		transaction.PaymentToken = uuid.FromStringOrNil(res.Token)
		// u.midtransGateway.
	}

	transaction.UpdatedAt = now
	transaction.UpdatedBy = userId

	err = u.useTransaction.UpdateHeader(ctx, Claims, transaction.Id, transaction)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Receive implements ipayment.Usecase
func (u *usePayment) Receive(ctx context.Context, request *models.MidtransNotification) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		logger = logging.Logger{}
	)

	logger.Info(request)

	switch request.TransactionStatus {
	case "capture":
		if request.PaymentType == "credit_card" && request.FraudStatus == "accept" {
			//
			err := u.PayTransaction(ctx, request)
			if err != nil {
				return err
			}
		} else {

			logger.Warn("transaction captured, potentially fraud")
			return nil
		}
	case "settlement", "deny", "cancel", "expire", "pending":
		err := u.PayTransaction(ctx, request)
		if err != nil {
			return err
		}
	default:
		logger.Warn("payment status type is unidentified")
		return nil
	}

	//save request to log
	var reqLog = models.MidtransNotificationLog{}
	err := mapstructure.Decode(request, &reqLog)
	if err != nil {
		logger.Error("error mapp decode request midtrans ", err)
	}

	reqLog.VaNumbers = fmt.Sprintf("%#v", request.VaNumbers)
	reqLog.PaymentAmounts = fmt.Sprintf("%#v", request.PaymentAmounts)
	err = u.repoPaymentNotifLog.Create(ctx, &reqLog)
	if err != nil {
		logger.Error("error save request midtrans ", err)
	}

	return nil
}

func (u *usePayment) PayTransaction(ctx context.Context, req *models.MidtransNotification) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		logger = logging.Logger{}
	)

	trx, err := u.repoTransaction.GetDataBy(ctx, "transaction_code", req.OrderID) //u.useTransaction.GetById(ctx, claim, req[0])
	if err != nil {
		logger.Error("payment.notification get transaction ", err)
		return err
	}

	if trx.StatusPayment != models.STATUS_PAYMENTSUCCESS {
		trx.StatusPayment = paymentStatus(req.TransactionStatus)
		trx.Description = req.PaymentType
		trx.PaymentStatusDesc = req.TransactionStatus
		trx.PaymentId = uuid.FromStringOrNil(req.TransactionID)
		err = u.repoTransaction.Update(ctx, trx.Id, trx)
		if err != nil {
			return err
		}
	}

	return nil
}

func paymentStatus(status string) models.StatusPayment {
	switch status {
	case "capture", "settlement":
		return models.STATUS_PAYMENTSUCCESS
	case "deny", "cancel":
		return models.STATUS_FAILED
	case "pending":
		return models.STATUS_WAITINGPAYMENT
	case "expire":
		return models.STATUS_EXPIRED
		// case "pending":
	}
	return 0
}

// Status implements ipayment.Usecase
func (u *usePayment) Status(ctx context.Context, Claims util.Claims, trxCode string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		logger = logging.Logger{}
	)
	data, err := u.coreGateway.CheckTransaction(trxCode)
	if err != nil {
		logger.Error("Failed check transaction midtrans ", err)
		return nil, err
	}
	return data, nil

}

// send notif
