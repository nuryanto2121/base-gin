package imidtrans

import (
	"github.com/midtrans/midtrans-go/coreapi"
)

type Repository interface {
	CheckTransaction(param string) (*coreapi.TransactionStatusResponse, error)
}
