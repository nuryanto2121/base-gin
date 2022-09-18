package itransaction

import (
	"app/models"
	"app/pkg/util"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetDataBy(ctx context.Context, key, value string) (result *models.Transaction, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result []*models.TransactionList, err error)
	GetListTicketUser(ctx context.Context, queryparam models.ParamList) (result []*models.TransactionResponse, err error)
	Create(ctx context.Context, data *models.Transaction) (err error)
	Update(ctx context.Context, ID uuid.UUID, data interface{}) (err error)
	Delete(ctx context.Context, ID uuid.UUID) (err error)
	Count(ctx context.Context, queryparam models.ParamList) (result int64, err error)
	CountUserList(ctx context.Context, queryparam models.ParamList) (result int64, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, transactionId string) (result *models.TransactionResponse, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	GetListTicketUser(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.TransactionForm) (result *models.TransactionResponse, err error)
	Payment(ctx context.Context, Claims util.Claims, data *models.TransactionPaymentForm) (err error)
	Update(ctx context.Context, Claims util.Claims, ID uuid.UUID, data *models.TransactionForm) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID uuid.UUID) (err error)
}
