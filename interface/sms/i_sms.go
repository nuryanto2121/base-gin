package isms

import (
	"app/models"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	Send(ctx context.Context, to, message, smsType string) (interface{}, error)
	Create(ctx context.Context, data *models.SmsLog) (err error)
}

type Usecase interface {
	Send(ctx context.Context, userId uuid.UUID, to, message, smsType string) error
}
