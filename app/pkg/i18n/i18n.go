package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"sync"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/gotnospirit/messageformat"
)

var Locales embed.FS

var localeToPlurals = map[string]string{
	"en": "en",
	"pt-BR": "pt",
}

type localeData struct {
	file   map[string]string
	parser *messageformat.Parser
}

var cache = make(map[string]localeData)
var mu sync.RWMutex

func getLocaleData(locale string) localeData {
	if item, ok := cache[locale]; ok {
		return item
	}

	mu.Lock()
	defer mu.Unlock()

	if item, ok := cache[locale]; ok {
		return item
	}

	content, err := Locales.ReadFile("locale/" + locale + ".json")
	if err != nil {
		panic(err)
	}

	var file map[string]string
	err = json.Unmarshal(content, &file)
	if err != nil {
		panic(err)
	}

	parser, err := messageformat.NewWithCulture(localeToPlurals[locale])
	if err != nil {
		panic(err)
	}

	data := localeData{file, parser}

	if env.IsProduction() {
		cache[locale] = data
	}

	return data
}

func getMessage(locale, key string) (string, *messageformat.Parser) {
	localeData := getLocaleData(locale)
	if str, ok := localeData.file[key]; ok {
		return str, localeData.parser
	}

	enData := getLocaleData("en")
	return enData.file[key], enData.parser
}

func T(ctx context.Context, key string, params ...map[string]interface{}) string {
	locale, ok := ctx.Value(app.LocaleCtxKey).(string)
	if !ok {
		locale = env.Config.Locale
	}

	msg, parser := getMessage(locale, key)
	if len(params) == 0 {
		return msg
	}

	parsedMsg, err := parser.Parse(msg)
	if err != nil {
		panic(err)
	}

	str, err := parsedMsg.FormatMap(params[0])
	if err != nil {
		panic(err)
	}

	return str
}