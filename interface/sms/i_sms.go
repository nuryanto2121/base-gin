package isms

import (
	"app/models"
	"context"
)

type Repository interface {
	Send(ctx context.Context, to, message, smsType string) error
	Create(ctx context.Context, data *models.SmsLog) (err error)
}

type Usecase interface {
	Send(ctx context.Context, to, message, smsType string) error
}
