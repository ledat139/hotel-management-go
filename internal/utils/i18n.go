package utils

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *goi18n.Bundle
var defaultLanguage = language.English

func InitI18n() {
	bundle = goi18n.NewBundle(defaultLanguage)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	pathEng := filepath.Join("internal", "locales", "en.json")
	pathVie := filepath.Join("internal", "locales", "vi.json")
	loadMessageFile(pathEng)
	loadMessageFile(pathVie)
	log.Println("i18n Bundle initialized. Default language:", defaultLanguage)
}

func loadMessageFile(path string) {
	_, err := bundle.LoadMessageFile(path)
	if err != nil {
		log.Printf("Warning: Could not load translation file ' %s': %v\n", path, err)
	} else {
		log.Printf("Successfully loaded translation file:%s\n", path)
	}
}

func T(c *gin.Context, messageID string) string {
	lang, exists := c.Get("lang")
	if !exists {
		lang = "en"
	}
	log.Println("Lang used:", lang)
	localizer := i18n.NewLocalizer(bundle, lang.(string))
	translated, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		log.Printf("Missing translation for: %s", messageID)
		return messageID
	}
	return translated
}
