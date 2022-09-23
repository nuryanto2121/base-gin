package midtrans

import (
	"app/pkg/setting"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// NewGateway creates new midtrans payment gateway
func NewGateway(creds setting.Midtrans) *Gateway {

	snapClient := &snap.Client{
		ServerKey: creds.SecretKey,
	}

	switch setting.ServerSetting.RunMode {
	case "prod":
		snapClient.Env = midtrans.Production
		snapClient.HttpClient = midtrans.GetHttpClient(midtrans.Production)
	default:
		snapClient.Env = midtrans.Sandbox
		snapClient.HttpClient = midtrans.GetHttpClient(midtrans.Sandbox)
	}

	gateway := Gateway{
		serverKey:     creds.SecretKey,
		SnapV2Gateway: snapClient,
	}

	return &gateway
}

// Gateway stores go-midtrans gateway and client
type Gateway struct {
	serverKey     string
	SnapV2Gateway SnapGateway
}

// NotificationValidationKey returns midtrans server key used for validating
// midtransa transaction status
func (g Gateway) NotificationValidationKey() string {
	return g.serverKey
}

type SnapGateway interface {
	CreateTransaction(req *snap.Request) (*snap.Response, *midtrans.Error)
}
