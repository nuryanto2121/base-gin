package multilanguage

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Language *i18n.Bundle

func Setup() {
	Language = i18n.NewBundle(language.English)
	Language.RegisterUnmarshalFunc("json", json.Unmarshal)
	Language.LoadMessageFile("pkg/multiLanguage/en.json")
	Language.LoadMessageFile("pkg/multiLanguage/id.json")
}
