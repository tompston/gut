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
	ex1 := gut.Convert(User{}, gut.Type{Name: "MyCustomInterface"})

	// Generate both the interface for the struct and
	// also a type which holds an array of interfaces.
	// + optionally you can also rename it.
	ex2 := gut.Convert(MyRandomStruct{},
		gut.Type{IsArray: true, ArrayTypeName: "ArrayOfMyRandomStructs"})

	// concat all of the interfaces together
	interfaces := fmt.Sprintln(ex1, ex2)

	if err := gut.GenerateTypescriptInterfaces(
		"./example.gen.ts", interfaces, gut.DefaultSettings); err != nil {
		fmt.Println(err)
	}
}

type User struct {
	ID        uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	Comments  `json:"comments,omitempty"`
}

type Comments []struct {
	ID    int    `json:"comment_id"`
	Value string `json:"value"`
}

type MyRandomStruct struct {
	MyFloat             float64
	MyInterface         interface{}
	Ex1                 map[string]map[string]string `json:"ex_1"`
	IntArray            []int                        `json:"int_array"`
	OptionalStringArray []string                     `json:"opt_str_array,omitempty"`
}
