package iauth

import (
	"context"

	"app/models"
	util "app/pkg/util"
)

type Usecase interface {
	LoginCms(ctx context.Context, dataLogin *models.LoginForm) (output interface{}, err error)
	LoginMobile(ctx context.Context, dataLogin *models.LoginForm) (output interface{}, err error)
	ForgotPassword(ctx context.Context, dataForgot *models.ForgotForm) (err error)
	ResetPassword(ctx context.Context, dataReset *models.ResetPasswd) (err error)
	Register(ctx context.Context, dataRegister models.RegisterForm) (err error)
	Verify(ctx context.Context, dataVeriry models.VerifyForm) (output interface{}, err error)
	VerifyForgot(ctx context.Context, dataVeriry models.VerifyForgotForm) (output interface{}, err error)
	Logout(ctx context.Context, claim util.Claims, Token string) (err error)
	CheckPhoneNo(ctx context.Context, PhoneNo string) (err error)
}

type AuthUsecase interface {
	GetUser(ctx context.Context, id int64) (result *models.Users, err error)
	Verify(ctx context.Context)
}
