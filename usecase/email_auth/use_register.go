package useemailauth

import (
	"strings"

	templateemail "app/pkg/email"
	util "app/pkg/utils"
)

type Register struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	GenerateNo string `json:"generate_no"`
}

func (R *Register) SendRegister() error {
	subjectEmail := "Verifikasi Code"

	err := util.SendEmail(R.Email, subjectEmail, getVerifyBody(R))
	if err != nil {
		return err
	}
	return nil
}

func getVerifyBody(R *Register) string {
	verifyHTML := templateemail.VerifyCode

	verifyHTML = strings.ReplaceAll(verifyHTML, `{Name}`, R.Name)
	verifyHTML = strings.ReplaceAll(verifyHTML, `{GenerateCode}`, R.GenerateNo)
	return verifyHTML
}
