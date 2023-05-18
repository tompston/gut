package main

import (
	"fmt"

	"github.com/tompston/gut"
)

// go run examples/05-generic-struct/main.go
func main() {

	ex1 := gut.Convert(GenericWithAnObject{})
	ex2 := gut.Convert(GenericWithAnArray{})
	ex3 := gut.Convert(GenericInsideGeneric{})
	ex4 := gut.Convert(GenericInsideGenericInsideGeneric{})

	if err := gut.Generate("./example.gen.ts", fmt.Sprintln(ex1, ex2, ex3, ex4)); err != nil {
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
type GenericInsideGenericInsideGeneric StructWithGeneric[GenericInsideGeneric]

/*
export type UuidType = string;
export type BigIntType = BigInt;
export type DateType = Date;

export interface GenericWithAnObject {
  some_field: string;
  areas: { [key: string]: any };
}

export interface GenericWithAnArray {
  some_field: string;
  areas: string[];
}

export interface GenericInsideGeneric {
  some_field: string;
  areas: {
    some_field: string;
    areas: { [key: string]: any };
  };
}

export interface GenericInsideGenericInsideGeneric {
  some_field: string;
  areas: {
    some_field: string;
    areas: {
      some_field: string;
      areas: { [key: string]: any };
    };
  };
}


*/
