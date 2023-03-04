package gut

import (
	"strings"
	"testing"
	"unicode"

	. "github.com/tompston/gut/types"
)

func TestTypescriptCodegen(t *testing.T) {
	type test struct {
		generated_interface string
		expected_interface  string
	}

	tests := []test{
		/* Tests on structs ( No settings ) */
		{
			generated_interface: ToTypescript(SimpleStruct{}),
			expected_interface: `
			export interface SimpleStruct {
				MyString: string
			}`,
		},
		{
			generated_interface: ToTypescript(SimpleStructWithJsonTags{}),
			expected_interface: `
			export interface SimpleStructWithJsonTags {
				my_str: string
				my_int: number
			}`,
		},
		{
			generated_interface: ToTypescript(SimpleStructWithTimeFields{}),
			expected_interface: `
			export interface SimpleStructWithTimeFields {
				MyString: string
				CreatedAt: DateType
				updated_at?: DateType
				deleted_at: DateType
			}`,
		},
		{
			// Check if time.Duration is correctly converted
			generated_interface: ToTypescript(StructWithTimeDurationField{}),
			expected_interface: `
			export interface StructWithTimeDurationField {
				SomeValue: string
				CurrentTime: BigIntType
			}`,
		},

		{
			generated_interface: ToTypescript(StructWithMultipleTypes{}),
			expected_interface: `
			export interface StructWithMultipleTypes {
				my_str: string
				my_int: number
				my_int_32: number
				ArrayOfStrings: string[]
				opt_arr_of_ints?: number[]
				MyInterface: any
			}`,
		},
		{
			generated_interface: ToTypescript(StructWithReference{}),
			expected_interface: `
			export interface StructWithReference {
				my_str: string
				MyInt: number
				ref: {
				  my_float: number
				  timestamp: number
				}
				opt_ref?: {
				  my_float: number
				  timestamp: number
				}
			}`,
		},
		{
			generated_interface: ToTypescript(StructWithArrayOfReferences{}),
			expected_interface: `
			export interface StructWithArrayOfReferences {
				arr_of_ref: {
				  my_float: number
				  timestamp: number
				}[]
			}`,
		},
		{
			generated_interface: ToTypescript(StructWithUnspecifiedStructName{}),
			expected_interface: `
			export interface StructWithUnspecifiedStructName {
				SomeValue: string
				ReferenceStruct: {
				  my_float: number
				  timestamp: number
				}
			}`,
		},
		{
			generated_interface: ToTypescript(StructWithMaps{}),
			expected_interface: `
			export interface StructWithMaps {
				MyStr: string
				my_opt_int?: BigIntType
				Ex1: { [key: string]: { [key: number]: string } }
				ex_2: { [key: string]: { [key: string]: string } }
				ex_3?: { [key: number]: { [key: string]: string } }
				ex_4?: { [key: string]: { [key: number]: any } }
			}`,
		},
		{
			generated_interface: ToTypescript(Employees{}),
			expected_interface: `
			export type EmployeesArray = Employees[]

			export interface Employees {
			  user_id: UuidType
			  Username: string
			  opt_surename?: string
			  RandomInterface: any
			  opt_interface?: any
			}`,
		},
		/* Tests on structs ( with settings ) */
		{
			generated_interface: ToTypescript(SimpleStruct{}, Interface{}),
			expected_interface: `
			export interface SimpleStruct {
				MyString: string
			}`,
		},
		{
			generated_interface: ToTypescript(SimpleStruct{}, Interface{Name: "_my_custom_name"}),
			expected_interface: `
			export interface _my_custom_name {
				MyString: string
			}`,
		},
		{
			generated_interface: ToTypescript(SimpleStruct{}, Interface{Name: "SingleRow", IsArray: true}),
			expected_interface: `
			export type SingleRowArray = SingleRow[]

			export interface SingleRow {
			  MyString: string
			}`,
		},
		{
			generated_interface: ToTypescript(SimpleStruct{}, Interface{Name: "SingleRow", IsArray: true, ArrayTypeName: "MyCustomNameForExportedArrayOfSimpleStructs"}),
			expected_interface: `
			export type MyCustomNameForExportedArrayOfSimpleStructs = SingleRow[]

			export interface SingleRow {
			  MyString: string
			}`,
		},
		{
			generated_interface: ToTypescript(Employees{}, Interface{Name: "EmployeeInterface"}),
			expected_interface: `
			export type EmployeeInterfaceArray = EmployeeInterface[]

			export interface EmployeeInterface {
			  user_id: UuidType
			  Username: string
			  opt_surename?: string
			  RandomInterface: any
			  opt_interface?: any
			}`,
		},
		{
			generated_interface: ToTypescript(Employees{}, Interface{Name: "EmployeeInterface", IsArray: false}),
			// If the Struct is array, but interface settings set it to false,
			// ignore that setting
			expected_interface: `
			export type EmployeeInterfaceArray = EmployeeInterface[]

			export interface EmployeeInterface {
			  user_id: UuidType
			  Username: string
			  opt_surename?: string
			  RandomInterface: any
			  opt_interface?: any
			}`,
		},
		{
			generated_interface: ToTypescript(Employees{}, Interface{Name: "EmployeeInterface", ArrayTypeName: "MyArrayOfEmployees"}),
			// Use custom array type name, if it's provided
			expected_interface: `
			export type MyArrayOfEmployees = EmployeeInterface[]

			export interface EmployeeInterface {
			  user_id: UuidType
			  Username: string
			  opt_surename?: string
			  RandomInterface: any
			  opt_interface?: any
			}`,
		},
	}

	for _, tc := range tests {
		gen_ts := stripSpaces(tc.generated_interface)
		exp_ts := stripSpaces(tc.expected_interface)
		if gen_ts != exp_ts {
			t.Fatalf("expected: %v\n, got: %v\n", tc.expected_interface, tc.generated_interface)
		}
	}
}

// util func for removig whitespaces, so that you can compare
// if the generated code matches the expected string.
func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
