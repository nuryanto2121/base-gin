package app

import (
	"fmt"

	multilanguage "app/pkg/multiLanguage"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Gin struct {
	C *gin.Context
}

// type Response struct {
// 	Code int         `json:"code"`
// 	Msg  string      `json:"msg"`
// 	Data interface{} `json:"data"`
// }
type Response struct {
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Errors interface{} `json:"errors"`
	Error  string      `json:"error"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode int, message string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Msg:  message,
		Data: data,
	})
	// return
}

// Response error setting gin.JSON
func (g *Gin) ResponseErrorMulti(httpCode int, message string, errors interface{}) {
	g.C.JSON(httpCode, Response{
		Errors: errors,
		Msg:    message,
	})
}

func (g *Gin) ResponseError(httpCode int, err error) {
	lang := g.C.Request.Header.Get("Language")
	loc := i18n.NewLocalizer(multilanguage.Language, lang)
	fmt.Println(lang)
	message, errLang := loc.Localize(&i18n.LocalizeConfig{
		MessageID: err.Error(),
	})
	if errLang != nil {
		g.C.JSON(httpCode, Response{
			Error: err.Error(),
			Msg:   err.Error(),
		})
	} else {
		g.C.JSON(httpCode, Response{
			Error: message,
			Msg:   message,
		})
	}

}
