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

	data, isNested := raw.(map[string]any)
	if !isNested {
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
