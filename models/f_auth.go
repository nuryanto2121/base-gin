package models

import "time"

//LoginForm :
type LoginForm struct {
	Account  string `json:"account" valid:"Required"`
	Password string `json:"pwd" valid:"Required"`
	// FirebaseToken string `json:"firebase_token,omitempty"`
}

// RegisterForm :
type RegisterForm struct {
	Name    string `json:"name" valid:"Required"`
	PhoneNo string `json:"phone_no"`
	// Email              string           `json:"email" valid:"Required;Email"`
	Password           string           `json:"pwd" valid:"Required;MinSize(6)"`
	ConfirmasiPassword string           `json:"confirm_pwd" valid:"Required;MinSize(6)"`
	Childs             []*RegisterChild `json:"childs"`
}

type RegisterChild struct {
	Name string    `json:"name"`
	DOB  time.Time `json:"dob"`
}

// ForgotForm :
type ForgotForm struct {
	Account string `json:"account" valid:"Required"`
}

// ResetPasswd :
type ResetPasswd struct {
	Account       string `json:"account" valid:"Required"`
	Passwd        string `json:"pwd" valid:"Required"`
	ConfirmPasswd string `json:"confirm_pwd" valid:"Required"`
}

type ResetPasswdMobile struct {
	AccessToken   string `json:"access_token" valid:"Required"`
	Passwd        string `json:"pwd" valid:"Required"`
	ConfirmPasswd string `json:"confirm_pwd" valid:"Required"`
}

type ResetPasswdProfile struct {
	OldPasswd     string `json:"old_pwd" valid:"Required"`
	Passwd        string `json:"pwd" valid:"Required"`
	ConfirmPasswd string `json:"confirm_pwd" valid:"Required"`
	AccessToken   string `json:"access_token" valid:"Required"`
}

type VerifyForm struct {
	Email       string `json:"email" valid:"Required;Email"`
	PhoneNo     string `json:"phone_no" valid:"Required"`
	AccessToken string `json:"access_token" valid:"Required"`
}

type VerifyForgotForm struct {
	Email       string `json:"email,omitempty"`
	PhoneNo     string `json:"phone_no" valid:"Required"`
	AccessToken string `json:"access_token" valid:"Required"`
	Otp         string `json:"otp" valid:"Required"`
}

type SosmedForm struct {
	Name        string `json:"name" valid:"Required"`
	Email       string `json:"email" valid:"Required;Email"`
	AccessToken string `json:"access_token" valid:"Required"`
}

type LoginAdminWebForm struct {
	Email    string `json:"email" valid:"Required;Email"`
	Password string `json:"password" valid:"Required"`
}
