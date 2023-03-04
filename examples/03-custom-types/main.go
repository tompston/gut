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
	ex1 := gut.ToTypescript(User{},
		gut.Interface{Name: "MyCustomInterface"})

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

/*

// running this would result in a file with the
// following content (post formatting)

// This is a custom comment in the file
export type UuidType = string
export type BigIntType = number
export type DateType = string


export interface MyCustomInterface {
  user_id: UuidType
  username: string
  created_at: DateType
}
*/
