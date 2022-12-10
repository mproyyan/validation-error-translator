package tl

import (
	"fmt"
	"os"
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

func TestAddTranslations(t *testing.T) {
	Load("en")

	tls := M{
		"test": "just test",
		"nested": M{
			"c1": "ok",
		},
	}

	AddTranslations(tls, false)

	test, found := translations["test"]
	assert.True(t, found)
	assert.Equal(t, "just test", test.(string))

	nested, found := translations["nested"]
	assert.True(t, found)
	assert.Equal(t, "ok", nested.(M)["c1"])
}

func TestReplaceTranslations(t *testing.T) {
	Load("en")

	tls := M{
		"required": "gak boleh kosong coy",
		"min": M{
			"numeric": "gak boleh lebih dari :p ya coyyy",
		},
		"lte": "jadi bukan nested",
		"semver": M{
			"nested": "jadi nested coy",
		},
	}

	AddTranslations(tls, true)

	required, ok := translations["required"]
	assert.True(t, ok)
	assert.Equal(t, "gak boleh kosong coy", required.(string))

	min, ok := translations["min"]
	assert.True(t, ok)
	assert.Equal(t, "gak boleh lebih dari :p ya coyyy", min.(M)["numeric"])

	lte, ok := translations["lte"]
	assert.True(t, ok)
	assert.Equal(t, "jadi bukan nested", lte.(string))

	semver, ok := translations["semver"]
	assert.True(t, ok)
	assert.Equal(t, "jadi nested coy", semver.(M)["nested"])
}

func TestAddTranslationsFailed(t *testing.T) {
	Load("en")

	tls := M{
		// "error": 78787, //cause of error
		"test": M{
			"vvvv": M{ // cause of error
				"ghdjf": "jfkg",
			},
		},
	}

	assert.Panics(t, func() {
		AddTranslations(tls, false)
	})

	_, falsy := translations["test"]
	assert.False(t, falsy)
}

func TestLoadFromFile(t *testing.T) {
	file, _ := os.Open("lang/en.json")
	defer file.Close()

	LoadFromFile(file)

	assert.NotNil(t, translations)
}
