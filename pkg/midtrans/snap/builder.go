package snap

import (
	"app/models"
	"app/pkg/setting"
	"fmt"
	"math"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func NewBuilder(inv *models.InvoiceRequest) *Builder {

	var callback *snap.Callbacks
	defaultRedirectUrl := setting.AppSetting.UrlSucessPayment //os.Getenv("INVOICE_SUCCESS_REDIRECT_URL")
	if defaultRedirectUrl != "" {
		callback = &snap.Callbacks{Finish: defaultRedirectUrl}
	}
	if inv.Callback.SuccessRedirectURL != "" {
		callback = &snap.Callbacks{
			Finish: inv.Callback.SuccessRedirectURL,
		}
	}

	srb := &Builder{
		req: &snap.Request{
			Items:     &[]midtrans.ItemDetails{},
			Callbacks: callback,
		},
	}

	return srb.
		setTransactionDetails(inv).
		setCustomerDetail(inv).
		setExpiration(inv).
		setItemDetails(inv).
		AddPaymentMethods(inv)
}

type Builder struct {
	req *snap.Request
}

func (b *Builder) setItemDetails(inv *models.InvoiceRequest) *Builder {
	var out []midtrans.ItemDetails

	for _, item := range inv.Items {

		name := item.Name
		if len(item.Name) > 50 {
			runes := []rune(name)
			name = string(runes[0:50])
		}

		out = append(out, midtrans.ItemDetails{
			ID:           item.ID,
			Name:         name,
			Price:        int64(item.Price),
			Qty:          int32(item.Qty),
			Category:     item.Category,
			MerchantName: item.MerchantName,
		})
	}

	b.req.Items = &out

	return b
}

func (b *Builder) setCustomerDetail(inv *models.InvoiceRequest) *Builder {
	b.req.CustomerDetail = &midtrans.CustomerDetails{
		FName: inv.Customer.Name,
		Email: inv.Customer.Email,
		Phone: inv.Customer.PhoneNumber,
		BillAddr: &midtrans.CustomerAddress{
			FName: inv.Customer.Name,
			Phone: inv.Customer.PhoneNumber,
		},
	}
	return b
}

func (b *Builder) setExpiration(inv *models.InvoiceRequest) *Builder {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	invDate := inv.TransactionDate.In(loc)
	fmt.Println(invDate)
	duration := inv.Duration //DueDate.Sub(inv.TransactionDate)
	b.req.Expiry = &snap.ExpiryDetails{
		StartTime: invDate.Format("2006-01-02 15:04:05 -0700"),
		Unit:      "minute",
		Duration:  int64(math.Round(duration.Minutes())),
	}
	fmt.Println(b.req.Expiry)
	return b
}

func (b *Builder) setTransactionDetails(inv *models.InvoiceRequest) *Builder {
	b.req.TransactionDetails = midtrans.TransactionDetails{
		OrderID:  inv.TransactionCode,
		GrossAmt: int64(inv.TotalAmount),
	}
	return b
}

func (b *Builder) AddPaymentMethods(inv *models.InvoiceRequest) *Builder {
	// b.req.EnabledPayments = append(b.req.EnabledPayments, snap.SnapPaymentType(inv.Payment))snap

	b.req.EnabledPayments = append(b.req.EnabledPayments, snap.PaymentTypeGopay)
	b.req.EnabledPayments = append(b.req.EnabledPayments, snap.PaymentTypeBankTransfer)
	b.req.EnabledPayments = append(b.req.EnabledPayments, snap.PaymentTypeBCAVA)
	b.req.EnabledPayments = append(b.req.EnabledPayments, snap.PaymentTypeIndomaret)
	b.req.EnabledPayments = append(b.req.EnabledPayments, snap.PaymentTypeAlfamart)

	return b
}

func (b *Builder) SetCreditCardDetail(d *snap.CreditCardDetails) *Builder {
	b.req.CreditCard = d
	return b
}

func (b *Builder) Build() (*snap.Request, error) {
	return b.req, nil
}
