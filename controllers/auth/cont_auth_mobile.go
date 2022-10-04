package contauth

import (
	iauth "app/interface/auth"
	"app/models"
	app "app/pkg/app"
	"app/pkg/logging"
	"app/pkg/middleware"
	tool "app/pkg/tools"
	util "app/pkg/util"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContAuthMobile struct {
	useAuth iauth.Usecase
}

func NewContAuthMobile(e *gin.Engine, useAuth iauth.Usecase) {
	cont := &ContAuthMobile{
		useAuth: useAuth,
	}

	v1 := e.Group("/v1")
	v1.Use(middleware.Versioning())
	v1.POST("/login", cont.LoginMobile)
	v1.POST("/register", cont.Register)
	v1.POST("/forgot-password", cont.ForgotPassword)
	v1.POST("/verify", cont.VerifyOtp)
	v1.POST("/change-password", cont.ChangePassword)

	v1.Use(middleware.Authorize())
	v1.POST("/logout", cont.Logout)

}

// Login :
// @Summary auth
// @Tags Auth Mobile
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.LoginForm true "this model set from firebase"
// @Success 200 {object} app.Response
// @Router /v1/login [post]
func (u *ContAuthMobile) LoginMobile(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.LoginForm{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	out, err := u.useAuth.LoginMobile(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", out)
}

// Register :
// @Summary Register
// @Tags Auth Mobile
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.RegisterForm true "Body with file zip"
// @Success 200 {object} app.Response
// @Router /v1/register [post]
func (u *ContAuthMobile) Register(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.RegisterForm{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := u.useAuth.Register(ctx, form)
	if err != nil {
		appE.ResponseError(http.StatusBadRequest, err)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// Logout :
// @Summary logout
// @Security ApiKeyAuth
// @Tags Auth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Success 200 {object} app.Response
// @Router /v1/cms/logout [post]
func (u *ContAuthMobile) Logout(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		appE = app.Gin{C: e} // wajib
	)

	claims, err := app.GetClaims(e)
	if err != nil {
		appE.ResponseError(http.StatusNetworkAuthenticationRequired, err)
		return
	}
	Token := e.Request.Header.Get("Authorization")
	err = u.useAuth.Logout(ctx, claims, Token)
	if err != nil {
		appE.Response(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return
	}

	appE.Response(http.StatusOK, "Ok", nil)
}

// ForgotPassword :
// @Summary Forgot Password : for generate otp and send to user then go to '/v1/verify'
// @Tags Auth Mobile
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.ForgotForm true "account fill with phone no set from verify forgot otp"
// @Success 200 {object} app.Response
// @Router /v1/forgot-password [post]
func (u *ContAuthMobile) ForgotPassword(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		// client sa_models.SaClient

		form = models.ForgotForm{}
	)

	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	rest, err := u.useAuth.ForgotPassword(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "", rest)
}

// VerifyOtp :
// @Summary Verify otp if success go to change-password
// @Tags Auth Mobile
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.VerifyForgotForm true "account fill with phone no set from verify forgot otp"
// @Success 200 {object} app.Response
// @Router /v1/verify [post]
func (u *ContAuthMobile) VerifyOtp(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		form   = models.VerifyForgotForm{}
	)

	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	rest, err := u.useAuth.VerifyForgot(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "", rest)
}

// ChangePassword :
// @Summary Change Password
// @Tags Auth Mobile
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.ResetPasswdMobile true "access token set from verify forgot otp"
// @Success 200 {object} app.Response
// @Router /v1/change-password [post]
func (u *ContAuthMobile) ChangePassword(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		// client sa_models.SaClient

		form = models.ResetPasswdMobile{}
	)

	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := u.useAuth.ResetPasswordMobile(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Please Login", nil)
}
