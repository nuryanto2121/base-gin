package models

//LoginForm :
type LoginForm struct {
	Account  string `json:"account" valid:"Required"`
	Password string `json:"pwd" valid:"Required"`
	// FirebaseToken string `json:"firebase_token,omitempty"`
}

// RegisterForm :
type RegisterForm struct {
	Name     string `json:"name" valid:"Required"`
	PhoneNo  string `json:"phone_no"`
	Email    string `json:"email" valid:"Required;Email"`
	Password string `json:"password" valid:"Required;MinSize(6)"`
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
	PhoneNo     string `json:"phone_no"`
	AccessToken string `json:"access_token"`
	Otp         string `json:"otp"`
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
