package tl

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var replacers = []*replacer{
	newReplacer(replaceFieldPlaceholder),
	newReplacer(replaceParamPlaceholder),
	newReplacer(replaceTagPlaceholder),
	newReplacer(replaceActualTagPlaceholder),
	newReplacer(replaceStructNamespacePlaceholder),
	newReplacer(replaceStructFieldPlaceholder),
	newReplacer(replaceNamespacePlaceholder),
}

func replaceFieldPlaceholder(raw string, f validator.FieldError) string {
	return strings.Replace(raw, ":f", f.Field(), -1)
}

func replaceParamPlaceholder(raw string, f validator.FieldError) string {
	return strings.Replace(raw, ":p", f.Param(), -1)
}

func replaceTagPlaceholder(raw string, f validator.FieldError) string {
	return strings.Replace(raw, ":t", f.Tag(), -1)
}

func replaceActualTagPlaceholder(raw string, f validator.FieldError) string {
	return strings.Replace(raw, ":at", f.ActualTag(), -1)
}

func replaceStructNamespacePlaceholder(raw string, f validator.FieldError) string {
	return strings.Replace(raw, ":sn", f.StructNamespace(), -1)
}

func replaceNamespacePlaceholder(raw string, f validator.FieldError) string {
	return strings.Replace(raw, ":n", f.Namespace(), -1)
}

func replaceStructFieldPlaceholder(raw string, f validator.FieldError) string {
	return strings.Replace(raw, ":sf", f.StructField(), -1)
}
