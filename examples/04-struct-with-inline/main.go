package main

import (
	"fmt"
	"time"

	"github.com/tompston/gut"
)

func main() {

	ex1 := gut.Convert(StructWithInlinedFields{})

	if err := gut.Generate(
		"./example.gen.ts", ex1, gut.Settings{
			FirstLine:  "// This is a custom comment in the file\n",
			DateType:   "string",
			UuidType:   "string",
			BigIntType: "number",
		}); err != nil {
		fmt.Println(err)
	}
}

type StructWithInlinedFields struct {
	MyEmbeddedStruct              `json:",inline"`
	NotEmbeddedStruct             `json:"not_embedded_struct"`
	CustomField                   int
	StructWhichHasEmbeddedStructs `json:"this_should_hold_start_end_and_updated_at"`
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
