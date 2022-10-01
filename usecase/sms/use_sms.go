package usesms

import (
	isms "app/interface/sms"
	"app/models"
	"app/pkg/logging"
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type useSms struct {
	repoSMS        isms.Repository
	contextTimeOut time.Duration
}

func NewSMS(a isms.Repository, timeout time.Duration) isms.Usecase {
	return &useSms{
		repoSMS:        a,
		contextTimeOut: timeout,
	}
}

// Send implements isms.Usecase
func (u *useSms) Send(ctx context.Context, userId uuid.UUID, to string, message string, smsType string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var logger = logging.Logger{}

	responseSMS, err := u.repoSMS.Send(ctx, to, message, smsType)
	if err != nil {
		//logging error
		if restError, ok := responseSMS.(models.TwillioResponseError); ok {
			u.repoSMS.Create(ctx, &models.SmsLog{
				ToUserId:       userId,
				PhoneNo:        to,
				Message:        message,
				MessageTwillio: fmt.Sprintf("error : %s | %s", restError.Message, restError.MoreInfo),
				Code:           int64(restError.Code),
			})
		}
		logger.Error("failed send sms ", err, to)
		return err
	}
	// else {
	//logging success
	if rest, ok := responseSMS.(models.TwillioResponse); ok {
		err = u.repoSMS.Create(ctx, &models.SmsLog{
			ToUserId:       userId,
			PhoneNo:        to,
			Message:        message,
			MessageTwillio: fmt.Sprintf("%s | %s", rest.Status, rest.URI),
			Code:           200,
		})
		if err != nil {
			return err
		}
	}
	// }

	return nil
}
