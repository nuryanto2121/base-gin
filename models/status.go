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
	STATUS_FAILED         StatusPayment = 1000004
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
	case STATUS_FAILED:
		result = "Payment Failed"

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
	PAYMENT_GOPAY   PaymentCode = 3000001
	PAYMENT_ATM     PaymentCode = 3000002
	PAYMENT_CASH    PaymentCode = 3000003
	PAYMENT_CASHIER PaymentCode = 3000004
	PAYMENT_BCA_VA  PaymentCode = 3000005
	PAYMENT_ONLINE  PaymentCode = 3000006
)

func (s PaymentCode) String() string {
	var result string = ""
	switch s {
	case PAYMENT_GOPAY:
		result = "gopay"
	case PAYMENT_ATM:
		result = "bank_transfer"
	case PAYMENT_CASH:
		result = "tunai"
	case PAYMENT_CASHIER:
		result = "cashier"
	case PAYMENT_BCA_VA:
		result = "bca_va"
	case PAYMENT_ONLINE:
		result = "Online Payment"

	}

	return result
}
