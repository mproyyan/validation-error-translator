package tl

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed lang/*.json
var translationFiles embed.FS

// translations hold all error translations after get loaded
// from lang folder, files are loaded based on what you choose,
// if you choose english translation then en.json will be loaded
// and if you choose an unavailable translation it will cause panic
var translations map[string]any

func Load(locale string) {
	switch locale {
	case "en":
		loadTranslation(locale)
	default:
		panic(fmt.Sprintf("translation for %s not available", locale))
	}
}

func loadTranslation(locale string) {
	path := fmt.Sprintf("lang/%s.json", locale)

	// load tranlation file from embedded file
	content, err := translationFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// mapping loaded translation file to translations
	err = json.Unmarshal(content, &translations)
	if err != nil {
		panic(err)
	}
}
