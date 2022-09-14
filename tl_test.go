package tl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
