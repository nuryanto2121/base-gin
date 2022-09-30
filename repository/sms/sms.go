package reposms

import (
	isms "app/interface/sms"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"
	"context"
	"encoding/json"
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
func (r *repoSms) Send(ctx context.Context, to string, message string, smsType string) error {

	var (
		// Set initial variables
		logger     = logging.Logger{}
		accountSid = setting.TwilioCredential.Sid
		authToken  = setting.TwilioCredential.Token
		from       = setting.TwilioCredential.From
		urlString  = fmt.Sprintf("%s/%s/Message.json", setting.TwilioCredential.Url, accountSid)
		//   urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

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
		return err
	}
	fmt.Println(resp.Status)
	var data map[string]interface{}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &data)
	if err == nil {
		fmt.Println(data["sid"])
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		fmt.Println(data)
		logger.Info("sms terkirim")
	} else {
		fmt.Println(data)
		logger.Info("sms gagal")
	}

	return nil
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
