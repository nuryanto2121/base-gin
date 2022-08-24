package models

type AuditLogsList struct {
	AddAuditLogs
	OutletName string `json:"outlet_name"`
	SkuName    string `json:"sku_name"`
}
