package ipayment

import (
	"app/models"
	"app/pkg/util"
	"context"
)

type Usecase interface {
	Receive(ctx context.Context, request *models.MidtransNotification) error
	Payment(ctx context.Context, Claims util.Claims, data *models.TransactionPaymentForm) (result *models.MidtransResponse, err error)
	Status(ctx context.Context, Claims util.Claims, trxCode string) (interface{}, error)
}
