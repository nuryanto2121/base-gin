package repomidtrans

import (
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	imidtrans "app/interface/midtrans"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

// repoGateway stores go-midtrans gateway and client
type repoGateway struct {
	serverKey string
	// CoreV2Gateway CoreGateway
	Client midtrans.HttpClient
	url    string
}

// NewRepoGateway creates new midtrans payment gateway
func NewRepoGateway(creds setting.Midtrans) imidtrans.Repository {

	var url string = ""

	coreClient := &coreapi.Client{
		ServerKey: creds.SecretKey,
	}

	switch setting.ServerSetting.RunMode {
	case "prod":
		coreClient.Env = midtrans.Production
		coreClient.HttpClient = midtrans.GetHttpClient(midtrans.Production)
		url = `https://app.midtrans.com`
	default:
		coreClient.Env = midtrans.Sandbox
		coreClient.HttpClient = midtrans.GetHttpClient(midtrans.Sandbox)
		url = `https://app.sandbox.midtrans.com`
	}
	coreClient.Env.BaseUrl()

	gateway := repoGateway{
		serverKey: creds.SecretKey,
		// CoreV2Gateway: coreClient,
		Client: coreClient.HttpClient,
		url:    url,
	}

	return &gateway
}

// NotificationValidationKey returns midtrans server key used for validating
// midtransa transaction status
func (g repoGateway) NotificationValidationKey() string {
	return g.serverKey
}

func (c repoGateway) CheckTransaction(param string) (*coreapi.TransactionStatusResponse, error) {

	var (
		result = &coreapi.TransactionStatusResponse{}
		logger = logging.Logger{}
	)
	url := fmt.Sprintf("%s/snap/v1/transactions/%s/status", c.url, param)
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		logger.Error("error get transaction status midtrans ", err)
		fmt.Print(err.Error())
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("error read body transaction status midtrans ", err)
		return nil, err
	}
	fmt.Println(string(responseData))
	if response.StatusCode != http.StatusOK || response.StatusCode != http.StatusCreated {
		logger.Error("error status header not ok ")
		return nil, models.ErrPaymentTokenExpired //status.Code()//errors.New("bad parameter")
	}

	if err := json.Unmarshal(responseData, &result); err != nil {
		// panic(err)
		return nil, err
	}
	fmt.Println(result)
	return result, nil
}
