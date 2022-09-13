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

	}

	return result
}
