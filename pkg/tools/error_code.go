package tool

import (
	"net/http"

	"app/models"

	"github.com/sirupsen/logrus"
)

// GetStatusCode :
func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound, models.ErrEmailNotFound, models.ErrOtpNotFound, models.ErrEmailNotFound, models.ErrVersioningNotFound, models.ErrAccountNotFound:
		return http.StatusNotFound
	case models.ErrConflict, models.ErrAccountConflict, models.ErrAccountAlreadyExist:
		return http.StatusConflict
	case models.ErrUnauthorized, models.ErrInvalidLogin, models.ErrInvalidPassword:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
