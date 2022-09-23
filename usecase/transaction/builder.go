package usetransaction

import (
	"app/models"
	"app/pkg/midtrans/snap"
	"context"
	"fmt"
	"time"
)

func (u *useTransaction) BuildMidtrans(ctx context.Context, trx *models.Transaction) (*snap.Builder, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	//get data user
	User, err := u.repoCustomer.GetDataBy(ctx, "id", trx.CustomerId.String())
	if err != nil {
		return nil, err
	}

	trxDetail, err := u.repoTransDetail.GetList(ctx, models.ParamList{
		Page:       1,
		PerPage:    50,
		InitSearch: fmt.Sprintf("transaction_id = '%s'", trx.Id),
	})
	if err != nil {
		return nil, err
	}

	category := "gelang"
	items := []models.Item{}
	for _, val := range trxDetail {
		if !val.IsChildren {
			category = "non gelang"
		}
		item := models.Item{
			ID:           val.Id.String(),
			Name:         val.SkuName,
			Category:     category,
			MerchantName: "Playtopia",
			Description:  "",
			Qty:          int(val.ProductQty),
			Price:        val.Price,
			Currency:     "IDR",
		}

		items = append(items, item)
	}

	invRequest := &models.InvoiceRequest{
		TransactionCode: trx.TransactionCode,
		TransactionDate: trx.TransactionDate,
		Duration:        13 * time.Minute,
		TotalAmount:     trx.TotalAmount,
		Payment:         models.PaymentType(trx.PaymentCode.String()),
		Customer: models.Customer{
			Name:        User.Name,
			Email:       "",
			PhoneNumber: User.PhoneNo,
		},
		Items:    items,
		Callback: &models.Callback{},
	}
	// response, err :=
	return snap.NewBuilder(invRequest), nil
}
