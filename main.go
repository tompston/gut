package gut

import (
	"bytes"
	"fmt"
	"os"
	r "reflect"
	"regexp"
	"strings"
	"time"
)

// Settings struct is passed to the Generate function,
// in order to specify what types you want to use for the emitted typescript
// interface fields, when the structs include the following types:
//   - time.Time
//   - uuid.UUID
//   - Int64 / Uint64
type Settings struct {
	// Optional first line in the generated file. Useful if you
	// want to write some custom comments.
	FirstLine string
	// type for the emitted time.Time values that
	// is used in structs. Could be either "string"
	// or "Date"
	DateType string
	// type for the emitted uuid.UUID types. Default should
	// be string, but can be changed if needed
	UuidType string
	// Not sure how js / ts handles int64 or uint64 types, so
	// Specify what type you want to use. Can be either
	// "number" or "BigInt"
	BigIntType string
}

// Type is an optional struct that can be passed to the Convert function, which modifies the generated typescript interfaces
type Type struct {
	// Optional custom name for the generated typescript interface. (Default = name of the struct )
	Name string
	// if set to true, will append a type that will hold an array of the generated interfaces. (Default = false)
	IsArray bool
	// optional name for the type that holds the array of interfaces (default = Name + "Array")
	ArrayTypeName string
}

// toTS converts the passed down type to the corresponding typescript interface type.
func toTS(typ r.Type, typeMap map[string]r.Type, inline ...bool) string {

	var isInline bool
	if len(inline) > 0 && inline[0] {
		isInline = true
	}

	switch typ.Kind() {

	case r.Struct:
		sb := strings.Builder{}

		if !isInline {
			sb.WriteString(" {\n")
		}

		if typ == r.TypeOf(time.Time{}) {
			return "DateType"
		}

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.PkgPath != "" { // Skip unexported fields
				continue
			}

			if hasInlineJsonTag(field) {
				sb.WriteString(fmt.Sprintf("%v\n", toTS(field.Type, typeMap, true)))
			} else {
				sb.WriteString(fmt.Sprintf("%v: %v\n", typescriptFieldname(field), toTS(field.Type, typeMap)))
			}
		}

		if !isInline {
			sb.WriteString("}")
		}
		return sb.String()

	case r.Slice:
		return fmt.Sprintf("%v[]", toTS(typ.Elem(), typeMap))

	/*
		This is commented out because uuid.UUID is converted
		into a number[] type in the typescript interfaces,
		not the expected UuidType.
	*/
	// case r.Array:
	// 	return fmt.Sprintf("%v[]", toTS(typ.Elem(), typeMap))

	case r.Map:
		return fmt.Sprintf("{[key: %v]: %v}", toTS(typ.Key(), typeMap), toTS(typ.Elem(), typeMap))

	case r.Ptr:
		return toTS(typ.Elem(), typeMap)

	default:
		if typ.Name() != "" {
			if _, ok := typeMap[typ.Name()]; !ok {
				typeMap[typ.Name()] = typ
			}
		}

		switch typ.Kind() {
		case r.String:
			return "string"
		case r.Bool:
			return "boolean"
		case
			r.Float32, r.Float64,
			r.Int, r.Int8, r.Int16, r.Int32,
			r.Uint, r.Uint8, r.Uint16, r.Uint32:
			return "number"
		case r.Int64, r.Uint64:
			return "BigIntType"
		default:
			if refType, ok := typeMap[typ.Name()]; ok {
				if refType.Name() == "UUID" {
					return "UuidType"
				}
				return refType.Name()
			}
			return "any"
		}
	}
}

func parseStruct(structType r.Type, structTypes map[string]r.Type, typeSettings ...Type) string {
	var buffer bytes.Buffer

	typeName := structType.Name()

	if len(typeSettings) == 1 {
		gutType := typeSettings[0]
		if gutType.Name == "" {
			gutType.Name = typeName
		}
		if !isValidTypeName(gutType.Name) {
			panic(fmt.Sprintf("Invalid typescript interface name was provided! %v", gutType.Name))
		}

		// set typeName to hold the custom value, if it was provided
		typeName = gutType.Name

		if gutType.IsArray {
			array_type_name := gutType.ArrayTypeName
			if array_type_name == "" {
				array_type_name = fmt.Sprintf("%sArray", typeName)
			}
			buffer.WriteString(fmt.Sprintf("export type %s = %s[] \n\n", array_type_name, typeName))
		}
	}

	// Start of the type
	buffer.WriteString(fmt.Sprintf("export interface %s {\n", typeName))
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if hasInlineJsonTag(field) {
			buffer.WriteString(fmt.Sprintf("  %s\n", toTS(field.Type, structTypes, true)))
		} else {
			buffer.WriteString(fmt.Sprintf("  %s: %s\n", typescriptFieldname(field), toTS(field.Type, structTypes)))
		}

	}

	buffer.WriteString("}\n\n")
	return buffer.String()
}

func hasInlineJsonTag(field r.StructField) bool {
	return strings.Contains(field.Tag.Get("json"), ",inline")
}

// Convert converts the passed in struct into a typescript interface
// and returns it as a string. The function also allows for the 2nd optional
// param, which is used to optionally define the settings of the generated
// typescript interface.
//
// Example
//
//	ex1 := gut.Convert(MyStruct{})
//	ex2 := gut.Convert(MyStruct{}, gut.Type{Name: "MyStructCustomName", IsArray : true})
func Convert(i interface{}, typeSettings ...Type) string {

	_type := make(map[string]r.Type)
	_typeof := r.TypeOf(i)

	if structIsArray(i) {
		// if the input struct is an array and the settings are present,
		// create a typescript interface with the settings if the
		// Interface.Name is present
		if len(typeSettings) == 1 {
			settings := typeSettings[0]
			settings.IsArray = true
			if settings.Name == "" {
				panic("The name for the array of structs cannot be empty!")
			}
			return parseStruct(_typeof.Elem(), _type, settings)
		}
		// else, if the interface is an array, but the settings are not present, set the IsArray setting to true.
		return parseStruct(_typeof.Elem(), _type, Type{IsArray: true, Name: _typeof.Name()})
	}
	// if the input struct is not an array and the settings are present
	if len(typeSettings) == 1 {
		return parseStruct(_typeof, _type, typeSettings[0])
	}
	// If the settings array is not present and the struct is not an array
	return parseStruct(_typeof, _type)
}

/* convert the field name into a valid value, based on the json tags */
func typescriptFieldname(field r.StructField) string {
	json_tag := field.Tag.Get("json")
	tag_name := strings.Split(json_tag, ",")[0]
	if json_tag == "" {
		return field.Name
	} else if strings.Contains(json_tag, ",omitempty") {
		return fmt.Sprintf("%v? ", tag_name)
	} else {
		return fmt.Sprintf("%v ", tag_name)
	}
}

// Default settings for the generated typescript file. Be free to create a custom Settings struct if needed.
var defaultSettings = Settings{
	DateType:   "Date",
	UuidType:   "string",
	BigIntType: "BigInt",
}

// Generate function creates a file  and saves the passed in content to it + appends
// the settings types at the start. If the 3rd optional param is present, it will be used
// to override the default settings.
func Generate(filename string, content string, settings ...Settings) error {
	// if the settings are not present, use the default settings, else use the provided settings
	var s Settings
	if len(settings) == 1 {
		s = settings[0]
	} else {
		s = defaultSettings
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	head_and_interfaces := fmt.Sprintln(createHeader(s), content)
	_, err = file.Write([]byte(head_and_interfaces))
	if err != nil {
		return err
	}

	fmt.Println("\033[32m * CREATED\033[0m ", filename)
	return nil
}

func structIsArray(v interface{}) bool {
	rv := r.ValueOf(v)
	if rv.Kind() != r.Slice {
		return false
	}
	if rv.Type().Elem().Kind() != r.Struct {
		return false
	}
	return true
}

// Thanks chatGPT
func isValidTypeName(name string) bool {
	// Define a list of TypeScript reserved words
	reservedWords := []string{
		"abstract", "as", "async", "await", "break", "case", "catch", "class", "const", "continue",
		"debugger", "declare", "default", "delete", "do", "else", "enum", "export", "extends", "false",
		"finally", "for", "from", "function", "get", "if", "implements", "import", "in", "infer", "instanceof",
		"interface", "is", "keyof", "let", "module", "namespace", "never", "new", "null", "number", "object",
		"package", "private", "protected", "public", "readonly", "require", "return", "set", "static", "string",
		"super", "switch", "symbol", "this", "throw", "true", "try", "type", "typeof", "unique", "unknown", "var",
		"void", "while", "with", "yield",
	}

	// Define a regex pattern for a valid TypeScript identifier
	validNamePattern := regexp.MustCompile(`^[a-zA-Z_$][0-9a-zA-Z_$]*$`)

	// Check if the name matches the valid name pattern
	if !validNamePattern.MatchString(name) {
		return false
	}

	// Check if the name is a reserved word
	for _, word := range reservedWords {
		if strings.ToLower(name) == word {
			return false
		}
	}

	return true
}

func createHeader(s Settings) string {
	sb := strings.Builder{}

	// Append first line if exists
	if s.FirstLine != "" {
		sb.WriteString(s.FirstLine)
	}

	// Check for uuid type
	if s.UuidType == "" {
		sb.WriteString(fmt.Sprintf("export type UuidType = %s\n", "string"))
	} else {
		sb.WriteString(fmt.Sprintf("export type UuidType = %s\n", s.UuidType))
	}

	// Check for int64/uint64 types
	if s.BigIntType == "" {
		sb.WriteString(fmt.Sprintf("export type BigIntType = %s\n", "BigInt"))
	} else {
		sb.WriteString(fmt.Sprintf("export type BigIntType = %s\n", s.BigIntType))
	}

	// Check for time.Time type
	if s.DateType == "" {
		sb.WriteString(fmt.Sprintf("export type DateType = %s\n", "Date"))
	} else {
		sb.WriteString(fmt.Sprintf("export type DateType = %s\n", s.DateType))
	}

	// seperate type definitions from the generated interfaces
	sb.WriteString("\n\n")

	return sb.String()
}
