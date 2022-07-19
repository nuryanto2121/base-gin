package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/369-engineer/369backend/account/pkg/app"
	maria "gitlab.com/369-engineer/369backend/account/pkg/mariadb"
	version "gitlab.com/369-engineer/369backend/account/pkg/middleware/versioning"
	"gitlab.com/369-engineer/369backend/account/pkg/redisdb"
	"gitlab.com/369-engineer/369backend/account/pkg/setting"
	util "gitlab.com/369-engineer/369backend/account/pkg/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Versioning() gin.HandlerFunc {
	return func(e *gin.Context) {
		var (
			DeviceType       = e.Request.Header.Get("Device-Type")
			Versi            = e.Request.Header.Get("Version")
			ctx              = e.Request.Context()
			Version    (int) = 0
			err        error
		)

		Version, err = strconv.Atoi(Versi)
		if err != nil {
			resp := app.Response{
				Msg:  fmt.Sprintf("Versioning : %v", err),
				Data: nil,
			}
			e.JSON(http.StatusBadRequest, resp)
			return
		}

		if Version == 0 {
			resp := app.Response{
				Msg:   "Please Set Header Version",
				Data:  nil,
				Error: "Please Set Header Version",
			}
			e.JSON(http.StatusExpectationFailed, resp)
			return
		}

		verService := &version.AppVersion{
			DeviceType: DeviceType,
			Version:    Version,
		}
		dataVersion, err := verService.GetVersion(ctx, maria.Conn)
		if err != nil {
			resp := app.Response{
				Error: fmt.Sprintf("Versioning : %v", err),
				Data:  nil,
			}
			e.JSON(http.StatusBadRequest, resp)
			return
		}

		if dataVersion.MinVersion > Version {
			resp := app.Response{
				Msg:   "",
				Data:  dataVersion.Version,
				Error: "",
			}
			e.JSON(http.StatusHTTPVersionNotSupported, resp)
			return
		}

		// return next(e)
		e.Next()
	}
}

func Authorize() gin.HandlerFunc {
	return func(e *gin.Context) {
		var (
			code  = http.StatusOK
			msg   = ""
			data  interface{}
			token = "" //strings.Split(e.Request.Header.Get("Authorization"), "Bearer ")[1]
		)
		if e.Request.Header.Get("Authorization") == "" {
			token = ""
		} else {
			token = strings.Split(e.Request.Header.Get("Authorization"), "Bearer ")[1]
		}

		data = map[string]string{
			"token": token,
		}

		if token == "" {
			code = http.StatusNetworkAuthenticationRequired
			msg = "Auth Token Required"
		} else {
			// validasi JWT
			existToken := redisdb.GetSession(token)
			if existToken == "" {
				code = http.StatusUnauthorized
				msg = "Token Failed"
			} else {
				//Validasi Session
				claims, err := util.ParseToken(token)
				if err != nil {
					code = http.StatusUnauthorized
					switch err.(*jwt.ValidationError).Errors {
					case jwt.ValidationErrorExpired:
						msg = "Token Expired"
					default:
						msg = "Token Failed"
					}
				} else {
					var issuer = setting.AppSetting.Issuer
					valid := claims.VerifyIssuer(issuer, true)
					if !valid {
						code = http.StatusUnauthorized
						msg = "Issuer is not valid"
					}
					e.Set("claims", claims)
				}
				expiredDate := util.Int64ToTime(claims.StandardClaims.ExpiresAt)
				if expiredDate.Before(util.GetTimeNow()) {
					resp := app.Response{
						Msg:   "Token Expired",
						Data:  data,
						Error: "Token Expired",
					}
					e.AbortWithStatusJSON(http.StatusUnauthorized, resp)
					return
				}

			}
		}
		if code != http.StatusOK {
			resp := app.Response{
				Msg:   msg,
				Data:  data,
				Error: msg,
			}
			e.AbortWithStatusJSON(code, resp)
			return

		}
		e.Next()
	}
}

func CheckToken(token string) (*util.Claims, error) {

	var (
		msg    = ""
		claims = &util.Claims{}
		err    error
	)
	existToken := redisdb.GetSession(token)
	if existToken == "" {
		return nil, errors.New("inValid Token")
	} else {
		//Validasi Session
		claims, err = util.ParseToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				msg = "Token Expired"
			default:
				msg = "Token Failed"
			}
			return nil, errors.New(msg)
		} else {
			var issuer = setting.AppSetting.Issuer
			valid := claims.VerifyIssuer(issuer, true)
			if !valid {
				msg = "Issuer is not valid"
				return nil, errors.New(msg)
			}
		}
		expiredDate := util.Int64ToTime(claims.StandardClaims.ExpiresAt)
		if expiredDate.Before(util.GetTimeNow()) {
			msg = "Token Expired"
			return nil, errors.New(msg)
		}

	}
	return claims, nil

}
