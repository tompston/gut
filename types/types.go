package types

import (
	"time"

	"github.com/google/uuid"
)

type SimpleStruct struct {
	MyString string
}

type SimpleStructWithTimeFields struct {
	MyString          string
	CreatedAt         time.Time
	UpdatedAtOptional time.Time `json:"updated_at,omitempty"`
	DeletedAt         time.Time `json:"deleted_at"`
}

type SimpleStructWithJsonTags struct {
	MyString string `json:"my_str"`
	MyInt    int    `json:"my_int"`
}

type StructWithMultipleTypes struct {
	MyString            string   `json:"my_str"`
	MyInt               int      `json:"my_int"`
	MyInt32             int32    `json:"my_int_32"`
	ArrayOfStrings      []string ``
	OptionalArrayOfInts []int    `json:"opt_arr_of_ints,omitempty"`
	MyInterface         interface{}
}

type StructWithReference struct {
	MyString          string `json:"my_str"`
	MyInt             int
	Reference         ReferenceStruct `json:"ref"`
	OptionalReference ReferenceStruct `json:"opt_ref,omitempty"`
}

type StructWithArrayOfReferences struct {
	ArrayOfRefences []ReferenceStruct `json:"arr_of_ref"`
}

type StructWithUnspecifiedStructName struct {
	SomeValue string
	ReferenceStruct
}

type StructWithTimeDurationField struct {
	SomeValue   string
	CurrentTime time.Duration
}

type ReferenceStruct struct {
	MyFloat   float64 `json:"my_float"`
	Timestamp int     `json:"timestamp"`
}

type StructWithMaps struct {
	MyStr    string
	MyOptInt int64 `json:"my_opt_int,omitempty"`
	Ex1      map[string]map[int]string
	Ex2      map[string]map[string]string       `json:"ex_2"`
	Ex3      map[int]map[string]string          `json:"ex_3,omitempty"`
	Ex4      map[string]map[float32]interface{} `json:"ex_4,omitempty"`
}

type Employees []struct {
	ID                uuid.UUID `json:"user_id"`
	Username          string
	OptionalSurename  string `json:"opt_surename,omitempty"`
	RandomInterface   interface{}
	OptionalInterface interface{} `json:"opt_interface,omitempty"`
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

type StructWithGeneric[T any] struct {
	SomeField   string `json:"some_field"`
	GenericType T      `json:"areas" bson:"areas"`
}

type GenericWithAnObject StructWithGeneric[map[string]interface{}]
type GenericWithAnArray StructWithGeneric[[]string]
type GenericInsideGeneric StructWithGeneric[GenericWithAnObject]
type GenericInsideGenericInsideGeneric StructWithGeneric[GenericInsideGeneric]
