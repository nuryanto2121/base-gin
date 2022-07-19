package app

import (
	"fmt"
	"net/http"

	"app/models"
	multilanguage "app/pkg/multiLanguage"
	util "app/pkg/utils"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) (int, string) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, models.ErrBadParamInput.Error()
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, models.ErrInternalServerError.Error()
	}
	fmt.Printf("\n%v", valid.Errors)
	if !check {
		return http.StatusBadRequest, MarkErrors(valid.Errors)
	}

	return http.StatusOK, "ok"
}

// BindAndValidMulti binds and validates data
func BindAndValidMulti(c *gin.Context, form interface{}) (int, interface{}) {

	lang := c.Request.Header.Get("Language")
	loc := i18n.NewLocalizer(multilanguage.Language, lang)
	fmt.Println(lang)

	result := []map[string]interface{}{}
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, models.ErrBadParamInput.Error()
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, models.ErrInternalServerError.Error()
	}
	if !check {
		for _, err := range valid.Errors {
			val := ""
			if err.LimitValue != nil {
				val = fmt.Sprintf("%v", err.LimitValue)
			}
			message, errLang := loc.Localize(&i18n.LocalizeConfig{
				MessageID: "message_error_" + err.Name,
				TemplateData: map[string]interface{}{
					"Name":  err.Field,
					"Value": val,
				},
			})
			if errLang != nil {
				message = err.Message
			}

			errs := make(map[string]interface{})
			errs[err.Field] = message
			result = append(result, errs)
		}
		return http.StatusBadRequest, result
	}

	return http.StatusOK, "ok"
}

func ValidEmail(Email string) (interface{}, bool) {
	result := []map[string]interface{}{}
	if !util.CheckEmail(Email) {
		errs := make(map[string]interface{})
		errs["email"] = models.ErrFormatEmail.Error()
		result = append(result, errs)
		return result, false
	}
	return nil, true
}

func ValidPhoneNo(PhoneNo string) (interface{}, bool) {
	result := []map[string]interface{}{}
	if PhoneNo == "" {
		errs := make(map[string]interface{})
		errs["phone_no"] = models.ErrRequreidPhoneNo.Error()
		result = append(result, errs)
		return result, false
	}
	return nil, true
}

// GetClaims :
func GetClaims(c *gin.Context) (util.Claims, error) {
	var clm util.Claims
	// claims := c.GetStringMapString("claims")
	claimsx := c.Keys["claims"]

	err := mapstructure.Decode(claimsx, &clm)
	if err != nil {
		return clm, err
	}

	return clm, nil
}
