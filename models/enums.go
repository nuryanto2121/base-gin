package models

type StatusOrder int64

const (
	SUBMITTED StatusOrder = iota
	APPROVE
	REJECT
)

type StatusPayment int64

const (
	STATUS_WAITINGPAYMENT StatusPayment = 1000001
	STATUS_PAYMENTSUCCESS StatusPayment = 1000002
	STATUS_EXPIRED        StatusPayment = 1000003
)

func (s StatusPayment) String() string {
	var result string = ""
	switch s {
	case STATUS_WAITINGPAYMENT:
		result = "Waiting Payment"
	case STATUS_PAYMENTSUCCESS:
		result = "Payment Success"
	case STATUS_EXPIRED:
		result = "Payment Expired"

	}

	return result
}

type StatusTransaction int64

const (
	STATUS_ORDER    StatusTransaction = 2000001
	STATUS_CHECKIN  StatusTransaction = 2000002
	STATUS_CHECKOUT StatusTransaction = 2000003
	STATUS_DRAF     StatusTransaction = 2000004
)

func (s StatusTransaction) String() string {
	var result string = ""
	switch s {
	case STATUS_ORDER:
		result = "Order"
	case STATUS_CHECKIN:
		result = "Check In"
	case STATUS_CHECKOUT:
		result = "Check Out"
	case STATUS_DRAF:
		result = "Draf"

	}

	return result
}

type PaymentCode int64

const (
	PAYMENT_BCA     PaymentCode = 3000001
	PAYMENT_CC      PaymentCode = 3000002
	PAYMENT_CASH    PaymentCode = 3000003
	PAYMENT_CASHIER PaymentCode = 3000004
	PAYMENT_OTHER   PaymentCode = 3000005
)

func (s PaymentCode) String() string {
	var result string = ""
	switch s {
	case PAYMENT_BCA:
		result = "BCA"
	case PAYMENT_CC:
		result = "Kartu Kredit"
	case PAYMENT_CASH:
		result = "Tunai"
	case PAYMENT_CASHIER:
		result = "Cashier"
	case PAYMENT_OTHER:
		result = "Other"

	}

	return result
}
