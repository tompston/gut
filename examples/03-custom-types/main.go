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
	ex1 := gut.Convert(User{},
		gut.Type{Name: "MyCustomInterface"})

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

type User struct {
	ID        uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
