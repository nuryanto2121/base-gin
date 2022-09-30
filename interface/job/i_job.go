package ijob

import "context"

type Usecase interface {
	NotificationSms(ctx context.Context) error
}
