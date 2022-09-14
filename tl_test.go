package tl

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var validation = validator.New()

type sample struct {
	Name string   `validate:"required"`
	Num  int      `validate:"lte=100"`
	Data []string `validate:"max=1"`
}

func TestLoadTranlation(t *testing.T) {
	assert.Nil(t, translations)
	Load("en")

	// check if translation loaded successfully
	assert.NotNil(t, translations)
	fmt.Println(translations)
}

func TestLoadTranslationFailed(t *testing.T) {
	assert.Panics(t, func() {
		Load("nothing")
	})

	assert.Nil(t, translations)
}

func TestTranslate(t *testing.T) {
	Load("en")

	s := sample{
		Num:  900,
		Data: []string{"hgfjksd", "ghfshg"},
	}

	err := validation.Struct(s)
	if err != nil {
		ferr := err.(validator.ValidationErrors)
		for _, f := range ferr {
			fmt.Println(Translate(f))
		}
	}
}

func TestAddTranslation(t *testing.T) {
	Load("en")

	AddTranslation("test", "just test", false)

	tl, ok := translations["test"]
	assert.True(t, ok)
	assert.Equal(t, "just test", tl.(string))
}

func TestReplaceTranslation(t *testing.T) {
	Load("en")

	AddTranslation("required", "changed", true)

	tl, ok := translations["required"]
	assert.True(t, ok)
	assert.Equal(t, "changed", tl.(string))
}
