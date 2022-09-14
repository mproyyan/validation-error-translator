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
