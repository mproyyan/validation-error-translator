<h1 align="center">Validation Error Translator</h1>

## Intrroduction
Validation error translator is a package that translate [github.com/go-playground/validator](https://github.com/go-playground/validator) error into human readable error based on the validation tag. <br>
ex: tag `required` will be translated to `The field cannot be empty.`

## Instalation
Use go get.
```console
go get github.com/mproyyan/validation-error-translator
```

## Usage
```go
package main

import (
   "fmt"

   "github.com/go-playground/validator/v10"
   tl "github.com/mproyyan/validation-error-translator"
)

type Sample struct {
   Name  string   `validate:"required"`
   Num   int      `validate:"lte=100"`
}

func main() {
   // load translation
   tl.Load("en")

   validation := validator.New()
   s := Sample{
      Name: "",
      Num:  200,
   }

   err := validation.Struct(s)
   if err != nil {
      for _, field := range err.(validator.ValidationErrors) {
         errTl := tl.Translate(fiels)
         fmt.Prinln(errTl)
      }
   }
}
```
##### Result
```
The Name cannot be empty.
The Num must be less than or equal to 100.
```

## Load from file
This package support to load translation from your own file, make sure file is json format

```json
// indonesian translation id.json

{
   "required": "Field :f tidak boleh kosong.",
   "lte": {
      "item": "Field :f harus kurang dari atau sama dengan :p item.",
      "string": "Field :f harus kurang dari atau sama dengan :p karakter.",
      "numeric": "Field :f harus kurang dari atau sama dengan :p."
   }
}
```

```go
func main() {
   file, _ := os.Open("path-to-your-file/id.json")
   defer file.Close()

   tl.LoadFromFile(file)
}
```

## Placeholder
Raw string translation that contains placeholder will be replaced with actual data, data obtained from validator.FieldError interface, you can see [here](https://github.com/go-playground/validator/blob/master/errors.go#L83) for the detail.

#### Placeholder list
-  `:f` replaced with (validator.FieldError).Field()
-  `:p` replaced with (validator.FieldError).Param()
-  `:t` replaced with (validator.FieldError).Tag()
-  `:at` replaced with (validator.FieldError).ActualTag()
-  `:n` replaced with (validator.FieldError).Namespace()
-  `:sn` replaced with (validator.FieldError).StructNamespace()
-  `:sf` replaced with (validator.FieldError).StructField()

## Rules adding new translations
1. Translation must be string cannot be int nor array
   
   ```json
   {
      "required": 1342, // invalid
      "required": ["test", 123], // invalid
      "required": "The :f cannot be empty." // valid
   }
   ```
2. Value of nested translation must be string cannot be int, array nor other nested translation and must have these three child (item, string, numeric)

   ```json
   {
      // invalid
      "len": {
         "item": 123,
         "string": [123, "test"],
         "numeric": {
            "invalid": "invalid"
         }
      },

      // valid
      "len": {
         "item": "The :f must have :p items.",
         "string": "The :f must be :p characters in length.",
         "numeric": "The :f must be equal to :p."
      }
   }
   ```

## Adding translation
#### Adding or replace single translation

```go
// adding test to translations
tl.AddTranslation("test", "this is just test", false)

// replace required translation
tl.AddTranslation("required", "replaced", true)
```

#### Adding or replace batch of translations

```go
// you can use map[string]any or tl.M
// tl.M is shorthand for map[string]any
newTranslations := tl.M{
   "test": "just test",
   "test2": "another test",
}

// adding translations
tl.AddTranslations(newTranslations, false)

replaceTranslations := map[string]any{
   "required": "replaced",
   "len": tl.M{
      "item": "replaced",
      "string": "replaced",
      "numeric": "replaced",
   },
}

// replace translations
tl.AddTranslations(replaceTranslations, true)
```