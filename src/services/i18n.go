package services

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/xYurii/Bell/src/database/schemas"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle
var Languages = map[string]string{
	"pt-BR": "pt",
	"en-US": "en",
}

func Translate(key string, user *schemas.User, data ...interface{}) string {
	localizer := i18n.NewLocalizer(Bundle, user.Language)
	config := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key,
		},
	}

	if len(data) > 0 {
		config.TemplateData = data[0]
	}

	return localizer.MustLocalize(config)
}

func init() {
	Bundle = i18n.NewBundle(language.Portuguese)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, lang := range Languages {
		basePath := fmt.Sprintf("src/locales/%s", lang)

		if files, err := os.ReadDir(basePath); err == nil {
			for _, file := range files {
				path := fmt.Sprintf("%s/%s", basePath, file.Name())
				_, err = Bundle.LoadMessageFile(path)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
