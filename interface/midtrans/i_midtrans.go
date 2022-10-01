package imidtrans

import (
	"github.com/midtrans/midtrans-go/coreapi"
)

type Repository interface {
	CheckTransaction(paymentToken string) (*coreapi.TransactionStatusResponse, error)
}
