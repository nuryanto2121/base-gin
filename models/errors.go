package models

import (
	"errors"
)

var (
	ErrInternalServerError       = errors.New("internal_server_error") // ErrInternalServerError : will throw if any the Internal Server Error happen
	ErrQueryDbError              = errors.New("query_db_error")
	ErrNotFound                  = errors.New("data_not_found") // ErrNotFound : will throw if the requested item is not exists
	ErrOtpNotFound               = errors.New("otp_not_found")  //ErrOtpNotFound
	ErrExpiredOtp                = errors.New("expired_otp")
	ErrConflict                  = errors.New("conflict")          // ErrConflict : will throw if the current action already exists
	ErrAccountConflict           = errors.New("conflict_account")  // ErrConflict : will throw if the current action already exists
	ErrPhoneNoConflict           = errors.New("conflict_phone_no") // ErrConflict : will throw if the current action already exists
	ErrBadParamInput             = errors.New("bad_parameter")     // ErrBadParamInput : will throw if the given request-body or params is not valid
	ErrUnauthorized              = errors.New("unauthorized")
	ErrInvalidLogin              = errors.New("invalid_login")
	ErrAccountNotFound           = errors.New("account_not_found")
	ErrAccountNotActive          = errors.New("account_not_active")
	ErrAccountAlreadyExist       = errors.New("account_already_exist")
	ErrDataAlreadyExist          = errors.New("data_already_exist")
	ErrVersioningNotFound        = errors.New("versioning_not_found")
	ErrVersioningHeaderNotFound  = errors.New("versioning_header_not_found")
	ErrUpdateYourApp             = errors.New("update_your_app")
	ErrForceUpdateYourApp        = errors.New("force_update_your_app")
	ErrInvalidPassword           = errors.New("invalid_password")
	ErrInvalidOldPassword        = errors.New("invalid_old_password")
	ErrWrongPasswordConfirm      = errors.New("password_and_confirm_not_same")
	ErrNoUploadFile              = errors.New("file_upload_not_found")
	ErrNoEmailCantForgotPassword = errors.New("no_email_cant_forgot_password")
	ErrFirstAccountNeeded        = errors.New("first_account_needed")
	ErrEmailNotFound             = errors.New("email_not_found")
	ErrExpiredAccessToken        = errors.New("expired_access_token")
	ErrInvalidOTP                = errors.New("invalid_otp")
	ErrFormatEmail               = errors.New("format_email")
	ErrRequreidPhoneNo           = errors.New("required_phone_no")
	ErrRequestAccessToken        = errors.New("request_access_token")
	ErrPushNotificationFailed    = errors.New("push_notification_failed")
	ErrClaimsDecode              = errors.New("claims_decode")
	ErrNotPaidPayment            = errors.New("not_paid_payment")
	ErrStillHaveDraf             = errors.New("still_have_draf")
	ErrTransactionNotFound       = errors.New("transaction_not_found")
	ErrPaymentNeeded             = errors.New("payment_needed")
	ErrPaymentTokenExpired       = errors.New("payment_token_expired")
	ErrNoStatusOrder             = errors.New("no_status_order")
	ErrNoStatusCheckIn           = errors.New("no_status_check_in")
	ErrNoStatusCheckOut          = errors.New("no_status_check_out")
	ErrInventoryNotFound         = errors.New("inventory_not_found")
	ErrQtyExceedStock            = errors.New("qty_exceed_stock")
	ErrOvertime                  = errors.New("overtime")
	ErrNoMatchOutlet             = errors.New("no_match_outlet")
)
