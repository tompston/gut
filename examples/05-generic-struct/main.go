package main

import (
	"fmt"

	"github.com/tompston/gut"
)

func main() {

	ex1 := gut.Convert(GenericWithAnObject{})
	ex2 := gut.Convert(GenericWithAnArray{})
	ex3 := gut.Convert(GenericInsideGeneric{})

	if err := gut.Generate("./example.gen.ts", fmt.Sprintln(ex1, ex2, ex3)); err != nil {
		fmt.Println(err)
	}
}

type StructWithGeneric[T any] struct {
	SomeField   string `json:"some_field"`
	GenericType T      `json:"areas" bson:"areas"`
}

type GenericWithAnObject StructWithGeneric[map[string]interface{}]
type GenericWithAnArray StructWithGeneric[[]string]
type GenericInsideGeneric StructWithGeneric[GenericWithAnObject]
