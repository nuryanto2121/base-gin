package usesms

import (
	isms "app/interface/sms"
	"app/pkg/logging"
	"context"
	"time"
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
func (u *useSms) Send(ctx context.Context, to string, message string, smsType string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var logger = logging.Logger{}

	err := u.repoSMS.Send(ctx, to, message, smsType)
	if err != nil {
		logger.Error("failed send sms ", err, to)
		return err
	}

	return nil
}
