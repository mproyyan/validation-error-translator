package tl

import (
	"embed"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

//go:embed lang/*.json
var translationFiles embed.FS

// translations hold all error translations after get loaded
// from lang folder, files are loaded based on what you choose,
// if you choose english translation then en.json will be loaded
// and if you choose an unavailable translation it will cause panic
var translations map[string]any

// replace raw translation into actual translation
// list of placeholders [:f, :p, :t, :at, :sn, :n, :sf]
// :f --> (validator.FieldError).Field()
// :p --> (validator.FieldError).Param()
// :t --> (validator.FieldError).Tag()
// :at --> (validator.FieldError).ActualTag()
// :sn --> (validator.FieldError).StructNamespace()
// :n --> (validator.FieldError).Namespace()
// :sf --> (validator.FieldError).StructField()
var rl *replacerList

func Load(locale string) {
	registerReplacer()

	switch locale {
	case "en":
		loadTranslation(locale)
	default:
		panic(fmt.Sprintf("translation for %s not available", locale))
	}
}

func loadTranslation(locale string) {
	path := fmt.Sprintf("lang/%s.json", locale)

	// load translation file from embedded file
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

// translate error based on validator tag
// if translation for current tag not found, default error will be use
func Translate(f validator.FieldError) string {
	// search translation based on validator tag
	raw, found := translations[f.Tag()]
	if !found {
		// return default error, because translation for current tag not found
		return f.Error()
	}

	if tl, ok := raw.(string); ok {
		// replace raw string translation into actual translation
		return rl.replace(tl, f)
	}

	var data M
	switch raw.(type) {
	case map[string]any:
		data = (M)(raw.(map[string]any))
	case M:
		data = raw.(M)
	default:
		panic("translation invalid, make sure it was string or use M if you want to make nested translation.")
	}

	nestedTl := getNestedTag(data, f)
	return rl.replace(nestedTl, f)
}

// nested tag used to check the type of field that we are validating
// so the translation will be different based on field type
// field type string will use string translation
// field type map, slice, array will use item translation
// if field type is not what is mentioned above, translation will be use numeric
// make sure if you have nested translation, the nested must have child (item, string, numeric)
// and the value of nested translation must be string
func getNestedTag(data map[string]any, f validator.FieldError) string {
	var translation string
	var asserted bool
	var kind reflect.Kind

	kind = f.Kind()
	if kind == reflect.Ptr {
		kind = f.Type().Elem().Kind()
	}

	switch kind {
	case reflect.String:
		if _, found := data["string"]; !found {
			return f.Error()
		}

		translation, asserted = data["string"].(string)
		if !asserted {
			panic("value of nested translation must be string.")
		}
	case reflect.Array, reflect.Slice, reflect.Map:
		if _, found := data["item"]; !found {
			return f.Error()
		}

		translation, asserted = data["item"].(string)
		if !asserted {
			panic("value of nested translation must be string.")
		}
	default:
		if _, found := data["numeric"]; !found {
			return f.Error()
		}

		translation, asserted = data["numeric"].(string)
		if !asserted {
			panic("value of nested translation must be string.")
		}
	}

	return translation
}

// adding single translation
func AddTranslation(tag, translation string, override bool) {
	// first search translation
	_, exists := translations[tag]

	// if translation didn't exist, then just add
	if !exists {
		translations[tag] = translation
	} else if exists && override {
		// if translation exists and override is true
		// then replace old translation with new one
		translations[tag] = translation
	}
}

// adding or replace batch of translations
func AddTranslations(tls M, override bool) {
	for tag, tl := range tls {
		existsTl, exists := translations[tag]
		if !exists {
			err := isTranslationValid(tl)
			if err != nil {
				panic(err)
			}

			translations[tag] = tl
			continue
		}

		switch tl.(type) {
		case string:
			if override {
				err := isTranslationValid(tl)
				if err != nil {
					panic(err)
				}
				translations[tag] = tl
			}
		case M:
			result := tl.(M).filter(existsTl, override)
			err := isTranslationValid(result)
			if err != nil {
				panic(err)
			}

			translations[tag] = result
		default:
			panic("translation invalid, make sure it was string or use M type if you want to make nested translation.")
		}
	}
}

// translation must be string
// if you want to make or replace nested transtion you must use M type
// the value of nested translation must be string
func isTranslationValid(tl any) error {
	switch tl.(type) {
	case string:
		return nil
	case M, map[string]any:
		var ms M
		if _, ok := tl.(map[string]any); ok {
			ms = (M)(tl.(map[string]any))
		} else {
			ms = tl.(M)
		}

		for _, m := range ms {
			if _, ok := m.(string); !ok {
				return fmt.Errorf("value of nested translation must be string")
			}
		}

		return nil
	default:
		return fmt.Errorf("translation invalid, make sure it was string or use M type if you want to make nested translation")
	}
}
