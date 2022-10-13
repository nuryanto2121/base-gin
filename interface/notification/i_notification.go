package inotification

import "context"

type Usecase interface {
	NotificationSms(ctx context.Context)
}
