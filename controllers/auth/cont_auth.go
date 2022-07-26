package contauth

import (
	"context"
	"fmt"
	"net/http"

	iauth "app/interface/auth"
	"app/models"
	app "app/pkg/app"
	"app/pkg/logging"
	"app/pkg/middleware"
	tool "app/pkg/tools"
	util "app/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ContAuth struct {
	useAuth iauth.Usecase
}

func NewContAuth(e *gin.Engine, useAuth iauth.Usecase) {
	cont := &ContAuth{
		useAuth: useAuth,
	}

	e.GET("/v1/cms/health_check", cont.Health)
	r := e.Group("/v1/cms")
	// r.Use(middleware.Versioning())
	r.POST("/login", cont.Login)
	r.POST("/forgot", cont.ForgotPassword)
	r.POST("/change-password", cont.ChangePassword)

	L := e.Group("/v1/cms/logout")
	L.Use(middleware.Authorize())
	L.POST("", cont.Logout)

	v1 := e.Group("/v1")
	v1.POST("/login", cont.LoginSosmed)
	v1.POST("/register", cont.Register)
	v1.Use(middleware.Authorize())
	v1.POST("/logout", cont.Logout)

}

func (u *ContAuth) Health(e *gin.Context) {
	e.JSON(http.StatusOK, "success")
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
func (u *ContAuth) Logout(e *gin.Context) {
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

// Login :
// @Summary Login
// @Tags Auth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.LoginForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/login [post]
func (u *ContAuth) Login(e *gin.Context) {
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

	out, err := u.useAuth.LoginCms(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Ok", out)
}

// Login :
// @Summary auth
// @Tags Auth Mobile
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.SosmedForm true "this model set from firebase"
// @Success 200 {object} app.Response
// @Router /v1/login [post]
func (u *ContAuth) LoginSosmed(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = app.Gin{C: e}
		form   = models.SosmedForm{}
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

// ChangePassword :
// @Summary Change Password
// @Tags Auth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.ResetPasswd true "account set from verify forgot otp"
// @Success 200 {object} app.Response
// @Router /v1/cms/change-password [post]
func (u *ContAuth) ChangePassword(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib
		// client sa_models.SaClient

		form = models.ResetPasswd{}
	)

	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := u.useAuth.ResetPassword(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Please Login", nil)
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
func (u *ContAuth) Register(e *gin.Context) {
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

// ForgotPassword :
// @Summary Forgot Password
// @Tags Auth
// @Produce json
// @Param Device-Type header string true "Device Type"
// @Param Version header string true "Version Apps"
// @Param Language header string true "Language Apps"
// @Param req body models.ForgotForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} app.Response
// @Router /v1/cms/forgot [post]
func (u *ContAuth) ForgotPassword(e *gin.Context) {
	ctx := e.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	var (
		logger = logging.Logger{} // wajib
		appE   = app.Gin{C: e}    // wajib

		form = models.ForgotForm{}
	)
	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValidMulti(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		appE.ResponseErrorMulti(http.StatusBadRequest, "Bad Parameter", errMsg)
		return
	}

	err := u.useAuth.ForgotPassword(ctx, &form)
	if err != nil {
		appE.ResponseError(tool.GetStatusCode(err), err)
		return
	}

	appE.Response(http.StatusOK, "Check Your Email", nil)

}
