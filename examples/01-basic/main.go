package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tompston/gut"
)

func main() {

	// define which structs you want to convert
	ex1 := gut.Convert(User{})
	ex2 := gut.Convert(Comments{})
	ex3 := gut.Convert(MyRandomStruct{})

	// concat all of the interfaces together
	interfaces := fmt.Sprintln(ex1, ex2, ex3)

	if err := gut.Generate(
		"./example.gen.ts",
		interfaces); err != nil {
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
	MyFloat     float64
	MyInterface interface{}
	Ex1         map[string]map[string]string
	IntArray    []int
}
