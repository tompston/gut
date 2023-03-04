# gut

> Convert Golang structs to Typescript interfaces in ~300 loc

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
	ex1 := gut.ToTypescript(User{})
	ex2 := gut.ToTypescript(Comments{})
	ex3 := gut.ToTypescript(MyRandomStruct{})

	// concat all of the interfaces together
	interfaces := fmt.Sprintln(ex1, ex2, ex3)

	if err := gut.GenerateTypescriptInterfaces(
		"./example.gen.ts",
		interfaces, gut.DefaultSettings); err != nil {
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
export type UuidType = string
export type BigIntType = BigInt
export type DateType = Date


export interface User {
  user_id: UuidType
  username: string
  created_at: DateType
  comments?: {
    comment_id: number
    value: string
  }[]
}

export type CommentsArray = Comments[]

export interface Comments {
  comment_id: number
  value: string
}

export interface MyRandomStruct {
  MyFloat: number
  MyInterface: any
  Ex1: { [key: string]: { [key: string]: string } }
  IntArray: number[]
}
```

</td></tr>
</tbody></table>

### Why?

I wanted to solve the problem of converting golang structs to typescript interfaces, so that the backend and frontend clients could be in-sync, but I found that the current packages which try to do this had bugs and did not have the features I wanted. So I had to write it from scratch (with some help from ChatGPT)

### Features

  - Handle cases when the convertable struct is an array
  - handle `uuid.UUID` & `time.Time` conversion
  - Avoid duplicate interface names, by generating only one typescript interface which will hold all of the types that are present in the struct.
  
  - Flexible typescript type system for the converted `time.Time`, `uuid.UUID` and `int64` / `uint64` types.
  - optionally generate the type which holds an array of interfaces
  - Ability to optionally rename the generated typescript interface to a custom name
  - flexibility in exporting the converted go structs from packages  
      - `ToTypescript()` returns a string which holds the generated typescript interface, thus grouping the structs from multiple packages, concatenating them together and then saving them to a single file is quite trivial.
  - Keep the package simple.
    - `gut` exports only 2 funtions
      - `ToTypescript()` -> converts the struct into a ts string
      - `GenerateTypescriptInterfaces()` -> save the converted ts interfaces to a file + define the settings for the types

### Disclaimer

- The generated TS code is not formatted. Maybe I'll fix this later.
- Won't do any golang interface -> typescript conversion.
- There might be bugs.

### Credits

The initial draft of the code was generated by ChatGPT. After some heavy refactoring, this is the end result.

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
	ex1 := gut.ToTypescript(User{},
		gut.Interface{Name: "MyCustomInterface"})

	// Generate both the interface for the struct and
	// also a type which holds an array of interfaces.
	// + optionally you can also rename it.
	ex2 := gut.ToTypescript(MyRandomStruct{},
		gut.Interface{IsArray: true, ArrayTypeName: "ArrayOfMyRandomStructs"})

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
```

</td><td>

```ts
export type UuidType = string
export type BigIntType = BigInt
export type DateType = Date


export interface MyCustomInterface {
  user_id: UuidType
  username: string
  created_at: DateType
  comments?: {
    comment_id: number
    value: string
  }[]
}

export type ArrayOfMyRandomStructs = MyRandomStruct[]

export interface MyRandomStruct {
  MyFloat: number
  MyInterface: any
  ex_1: { [key: string]: { [key: string]: string } }
  int_array: number[]
  opt_str_array?: string[]
}
```

</td></tr>
</tbody></table>

### Example 3 - Change the predefined typesript type values

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
```

</td><td>

```ts
// This is a custom comment in the file
export type UuidType = string
export type BigIntType = number
export type DateType = string


export interface MyCustomInterface {
  user_id: UuidType
  username: string
  created_at: DateType
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