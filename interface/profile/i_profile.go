package iprofile

import (
	"context"

	"gitlab.com/369-engineer/369backend/account/models"
	util "gitlab.com/369-engineer/369backend/account/pkg/utils"
)

type Usecase interface {
	Update(ctx context.Context, claims util.Claims, data models.ProfileForm) (err error)
	GetDataBy(ctx context.Context, claims util.Claims, contact []models.ContactPhone) (result *models.ProfileResponse, err error)
	ResetPassword(ctx context.Context, claims util.Claims, dataReset *models.ResetPasswdProfile) (err error)
}
