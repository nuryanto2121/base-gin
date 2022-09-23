package usepayment

import (
	ipayment "app/interface/payment"
	itransaction "app/interface/transaction"
	"app/models"
	"app/pkg/logging"
	"app/pkg/util"
	"context"
	"time"

	_midtransGateway "app/pkg/midtrans"

	"github.com/midtrans/midtrans-go"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type usePayment struct {
	useTransaction  itransaction.Usecase
	repoTransaction itransaction.Repository
	contextTimeOut  time.Duration
	midtransGateway *_midtransGateway.Gateway
}

func NewUsePayment(a itransaction.Usecase, b *_midtransGateway.Gateway, c itransaction.Repository, timeout time.Duration) ipayment.Usecase {
	return &usePayment{
		useTransaction:  a,
		midtransGateway: b,
		repoTransaction: c,
		contextTimeOut:  timeout,
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

	transaction.UpdatedAt = now
	transaction.UpdatedBy = userId

	err = u.useTransaction.UpdateHeader(ctx, Claims, transaction.Id, transaction)
	if err != nil {
		return nil, err
	}

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
			err := u.PayTransaction(ctx, request.TransactionID, request.OrderID, request.TransactionStatus, request.StatusMessage)
			if err != nil {
				return err
			}
		} else {

			logger.Warn("transaction captured, potentially fraud")
			return nil
		}
	case "settlement", "deny", "cancel", "expire", "pending":
		err := u.PayTransaction(ctx, request.TransactionID, request.OrderID, request.TransactionStatus, request.StatusMessage)
		if err != nil {
			return err
		}
	default:
		logger.Warn("payment status type is unidentified")
		return nil
	}
	return nil
}

func (u *usePayment) PayTransaction(ctx context.Context, req ...string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		logger = logging.Logger{}
	)

	trx, err := u.repoTransaction.GetDataBy(ctx, "transaction_code", req[1]) //u.useTransaction.GetById(ctx, claim, req[0])
	if err != nil {
		logger.Error("payment.notification get transaction ", err)
		return err
	}

	trx.StatusPayment = paymentStatus(req[2])
	trx.Description = req[3]
	trx.PaymentId = uuid.FromStringOrNil(req[0])
	err = u.repoTransaction.Update(ctx, trx.Id, trx)
	if err != nil {
		return err
	}
	return nil
}

func paymentStatus(status string) models.StatusPayment {
	switch status {
	case "capture", "settlement":
		return models.STATUS_PAYMENTSUCCESS
	case "deny", "cancel", "pending":
		return models.STATUS_FAILED
	case "expire":
		return models.STATUS_EXPIRED
		// case "pending":
	}
	return 0
}

// send notif
