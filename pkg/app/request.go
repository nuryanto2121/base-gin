package app

import (
	"errors"

	"github.com/astaxie/beego/validation"

	"gitlab.com/369-engineer/369backend/account/pkg/logging"
)

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) string {
	result := ""
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
		result += ", " + err.Message
	}

	return result[2:]
}

func ErrorString(message string) error {
	return errors.New(message)
}
