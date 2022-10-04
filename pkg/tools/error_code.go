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
	case models.ErrNotFound, models.ErrEmailNotFound, models.ErrOtpNotFound, models.ErrEmailNotFound, models.ErrVersioningNotFound, models.ErrAccountNotFound, models.ErrTransactionNotFound:
		return http.StatusNotFound
	case models.ErrConflict, models.ErrAccountConflict, models.ErrAccountAlreadyExist, models.ErrDataAlreadyExist, models.ErrInventoryNotFound:
		return http.StatusConflict
	case models.ErrUnauthorized, models.ErrInvalidLogin, models.ErrInvalidPassword, models.ErrClaimsDecode, models.ErrPaymentTokenExpired, models.ErrExpiredAccessToken:
		return http.StatusUnauthorized
	case models.ErrPaymentNeeded:
		return http.StatusPaymentRequired
	case models.ErrNoStatusCheckIn, models.ErrNoStatusOrder, models.ErrNoStatusCheckOut, models.ErrQtyExceedStock:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
