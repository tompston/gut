# gut

> Convert Golang structs to Typescript interfaces

Handle

- json tags, like `omitempty` and `inline` appropriately
- `uuid.UUID` & `time.Time` conversion
- type generation for arrays
- setting custom names for the generated typescript interfaces
- types which use generics (see `examples/5-generic-struct`)

### Install

```bash
go get github.com/tompston/gut
```

### Example

<table>
<thead><tr><th>main.go</th><th>example.gen.ts</th></tr></thead>
<tbody>
<tr><td>

```go
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

	if err := gut.Generate("./example.gen.ts",interfaces); err != nil {
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
```

</td><td>

```ts
export type UuidType = string;
export type BigIntType = BigInt;
export type DateType = Date;

export interface User {
  user_id: UuidType;
  username: string;
  created_at: DateType;
  comments?: {
    comment_id: number;
    value: string;
  }[];
}

export type CommentsArray = Comments[];

export interface Comments {
  comment_id: number;
  value: string;
}

export interface MyRandomStruct {
  MyFloat: number;
  MyInterface: any;
  Ex1: { [key: string]: { [key: string]: string } };
  IntArray: number[];
}
```

</td></tr>
</tbody></table>

_To format the generated typescript code, you can run `deno fmt .`_

### Why?

I wanted to solve the problem of converting golang structs to typescript
interfaces, so that the backend and frontend clients could be in-sync, but I
found that the current packages which try to do this had bugs and did not have
the features I wanted. So I had to write it from scratch (with some help from
ChatGPT)

### Test

```bash
go test . -v -count=1
```

### Features

- Handle cases when the convertable struct is an array

- Handle json ",omitempty" tags
- Handle json ",inline" tags (embeded structs)
  - If a struct has a field with an json ",inline" tag, then the generated
    typescript interface will have all of the fields from the embeded struct
    inside of it. (see `examples/4-struct-with-inline`)
- handle `uuid.UUID` & `time.Time` conversion
- Avoid duplicate interface names, by generating only one typescript interface
  which will hold all of the types that are present in the struct.

- Flexible typescript type system for the converted `time.Time`, `uuid.UUID` and
  `int64` / `uint64` types.
- optionally generate the type which holds an array of interfaces
- Ability to optionally rename the generated typescript interface to a custom
  name
- flexibility in exporting the converted go structs from packages
  - `Convert()` returns a string which holds the generated typescript
    interface, thus grouping the structs from multiple packages, concatenating
    them together and then saving them to a single file is quite trivial.
- Keep the package simple.
  - `gut` exports only 2 funtions
    - `Convert()` -> converts the struct into a ts string
    - `Generate()` -> save the converted ts interfaces to a
      file + define the settings for the types

### Disclaimer

- The generated TS code is not formatted. Maybe I'll fix this later. The current way to format can be done by using `deno fmt .`
- There might be bugs.

### Credits

The initial draft of the code was generated by ChatGPT. After some heavy
refactoring, this is the end result.

### Example 2 - Modify interface names + array types

<table>
<thead><tr><th>main.go</th><th>example.gen.ts</th></tr></thead>
<tbody>
<tr><td>

```go
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

	// Generate both the interface for the struct and
	// also a type which holds an array of interfaces.
	// + optionally you can also rename it.
	ex2 := gut.Convert(MyRandomStruct{},
		gut.Type{IsArray: true, ArrayTypeName: "ArrayOfMyRandomStructs"})

	// concat all of the interfaces together
	interfaces := fmt.Sprintln(ex1, ex2)

	if err := gut.Generate("./example.gen.ts", interfaces); err != nil {
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
```

</td><td>

```ts
export type UuidType = string;
export type BigIntType = BigInt;
export type DateType = Date;

export interface MyCustomInterface {
  user_id: UuidType;
  username: string;
  created_at: DateType;
  comments?: {
    comment_id: number;
    value: string;
  }[];
}

export type ArrayOfMyRandomStructs = MyRandomStruct[];

export interface MyRandomStruct {
  MyFloat: number;
  MyInterface: any;
  ex_1: { [key: string]: { [key: string]: string } };
  int_array: number[];
  opt_str_array?: string[];
}
```

</td></tr>
</tbody></table>

### Example 3 - Change the predefined typescript type values

<table>
<thead><tr><th>main.go</th><th>example.gen.ts</th></tr></thead>
<tbody>
<tr><td>

```go
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
```

</td><td>

```ts
// This is a custom comment in the file
export type UuidType = string;
export type BigIntType = number;
export type DateType = string;

export interface MyCustomInterface {
  user_id: UuidType;
  username: string;
  created_at: DateType;
}
```

</td></tr>
</tbody></table>

<!--

## Creating a package

go mod init github.com/tompston/gut

git add .
git commit -m "gut: first release"
git tag v0.0.2
git push origin v0.0.2

 -->
