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
	util "app/pkg/util"

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
	r.Use(middleware.Versioning())
	r.POST("/login", cont.Login)
	r.POST("/logout", cont.Logout)
	r.Use(middleware.Authorize()).POST("/change-password", cont.ChangePassword)

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

// ChangePassword :
// @Summary Change Password
// @Security ApiKeyAuth
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
