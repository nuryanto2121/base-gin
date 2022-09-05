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

type StatusTransaction int64

const (
	STATUS_ORDER    StatusTransaction = 2000001
	STATUS_CHECKIN  StatusTransaction = 2000002
	STATUS_CHECKOUT StatusTransaction = 2000003
)
