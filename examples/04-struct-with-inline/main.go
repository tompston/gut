package main

import (
	"fmt"
	"time"

	"github.com/tompston/gut"
)

func main() {

	ex1 := gut.Convert(MyStruct{}, gut.Type{Name: "MyCustomInterface"})

	if err := gut.GenerateTypescriptInterfaces(
		"./example.gen.ts", ex1, gut.Settings{
			FirstLine:  "// This is a custom comment in the file\n",
			DateType:   "string",
			UuidType:   "string",
			BigIntType: "number",
		}); err != nil {
		fmt.Println(err)
	}
}

type MyStruct struct {
	MyEmbeddedStruct  `json:",inline"`
	NotEmbeddedStruct `json:"not_embedded_struct"`
	CustomField       int
	StructWhichHasEmbeddedStructs
}

type MyEmbeddedStruct struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotEmbeddedStruct struct {
	SomeRandomField string `json:"some_random_field"`
	SomeMoreStuff   map[string]interface{}
}

type StructWhichHasEmbeddedStructs struct {
	MyEmbeddedStruct `json:",inline"`
}
