package reposms

import (
	isms "app/interface/sms"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// type EnvSMS struct{

// }
type repoSms struct {
	key string
	db  db.DBGormDelegate
}

func NewRepoSMS(key string, Conn db.DBGormDelegate) isms.Repository {
	return &repoSms{
		key: key,
		db:  Conn,
	}
}

// Send implements isms.Repository
func (r *repoSms) Send(ctx context.Context, to string, message string, smsType string) (interface{}, error) {

	var (
		// Set initial variables
		logger     = logging.Logger{}
		accountSid = setting.TwilioCredential.Sid
		authToken  = setting.TwilioCredential.Token
		from       = setting.TwilioCredential.From
		urlString  = fmt.Sprintf("%s/%s/Messages.json", setting.TwilioCredential.Url, accountSid)
		respError  = models.TwillioResponseError{
			Code:   500,
			Status: 500,
		}
	)
	if setting.TwilioCredential.Mode == "debug" {
		to = setting.TwilioCredential.To
	}

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", to)
	v.Set("From", from)
	v.Set("Body", message)
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlString, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed send sms ", err)
		respError.Message = "error client do :" + err.Error()
		return respError, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data = models.TwillioResponse{}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			respError.Message = "error sms Unmarshal :" + err.Error()
			return respError, err
		}
		logger.Info("sms terkirim")
		return data, nil
	} else {
		fmt.Println(resp.Status)
		var data = models.TwillioResponseError{}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			return respError, err
		}

		fmt.Println(data)
		logger.Info("sms gagal")
		return data, errors.New(data.Message)
	}

}

// Create implements isms.Repository
func (r *repoSms) Create(ctx context.Context, data *models.SmsLog) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	conn := r.db.Get(ctx)
	query := conn.Create(data)

	err = query.Error
	if err != nil {
		logger.Error("repo sms Create ", err)
		return err
	}
	return nil
}
