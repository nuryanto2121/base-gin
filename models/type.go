package models

type PaymentType string

const (
	SourceCreditCard PaymentType = "credit_card"
	SourceBNIVA      PaymentType = "bni_va"
	SourcePermataVA  PaymentType = "permata_va"
	SourceBCAVA      PaymentType = "bca_va"
	SourceOtherVA    PaymentType = "other_va"
	SourceAlfamart   PaymentType = "alfamart"
	SourceGopay      PaymentType = "gopay"
	SourceAkulaku    PaymentType = "akulaku"
	SourceOvo        PaymentType = "ovo"
	SourceDana       PaymentType = "dana"
	SourceLinkAja    PaymentType = "linkaja"
	SourceShopeePay  PaymentType = "shopeepay"
	SourceQRIS       PaymentType = "qris"
	SourceBRIVA      PaymentType = "bri_va"
	SourceMandiriVA  PaymentType = "mandiri_va"
)

// Bank is a bank
type Bank string

const (
	BankBCA Bank = "bca"
	BankBNI Bank = "bni"
	BankBRI Bank = "bri"
)
