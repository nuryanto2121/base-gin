package iusersession

import (
	"app/models"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetByUser(ctx context.Context, UserId uuid.UUID) (result *models.UserSession, err error)
	GetByToken(ctx context.Context, Token string) (result *models.UserSession, err error)
	Create(ctx context.Context, data *models.UserSession) (err error)
	Delete(ctx context.Context, Token string) (err error)
}
