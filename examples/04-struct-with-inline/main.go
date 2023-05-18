package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tompston/gut"
)

func main() {

	// Insetad of generating an interface called User,
	// create one with a custom name
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

type User struct {
	ID        uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
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

type MyStruct struct {
	MyEmbeddedStruct  `json:",inline"`
	NotEmbeddedStruct `json:"not_embedded"`
	CustomField       int
}
