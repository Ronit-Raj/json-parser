# JSON Parser

This is a JSON parser written from scratch in Go . It aims to be fully compatible 
with [RFC8259](https://datatracker.ietf.org/doc/html/rfc8259) . 

## Usage Examples

### Step 1: Import the Package

```go
import (
    "fmt"
    "json_parser/parser"
)
```

### Step 2: Basic Types

#### Parsing a String
```go
json := `"hello world"`
var str string
if err := parser.Decode(json, &str); err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(str) // Output: hello world
```

#### Parsing a Number
```go
json := `42.5`
var num float64
if err := parser.Decode(json, &num); err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(num) // Output: 42.5
```

#### Parsing a Boolean
```go
json := `true`
var flag bool
if err := parser.Decode(json, &flag); err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(flag) // Output: true
```

#### Parsing Null
```go
json := `null`
var ptr *string
if err := parser.Decode(json, &ptr); err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(ptr == nil) // Output: true
```

### Step 3: Arrays

#### Simple Array
```go
json := `[1, 2, 3, 4, 5]`
var arr []any
if err := parser.Decode(json, &arr); err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(arr) // Output: [1 2 3 4 5]
```

#### Mixed Type Array
```go
json := `["hello", 42, true, null]`
var arr []any
if err := parser.Decode(json, &arr); err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(arr[0]) // Output: hello
fmt.Println(arr[1]) // Output: 42
fmt.Println(arr[2]) // Output: true
fmt.Println(arr[3]) // Output: <nil>
```

#### Empty Array
```go
json := `[]`
var arr []any
if err := parser.Decode(json, &arr); err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(len(arr)) // Output: 0
```

### Step 4: Objects

#### Simple Object
```go
json := `{"name": "Alice", "age": 30}`
var obj map[string]any
if err := parser.Decode(json, &obj); err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(obj["name"]) // Output: Alice
fmt.Println(obj["age"])  // Output: 30
```

#### Empty Object
```go
json := `{}`
var obj map[string]any
if err := parser.Decode(json, &obj); err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(len(obj)) // Output: 0
```

### Step 5: Nested Structures

#### Nested Objects
```go
json := `{
    "user": {
        "name": "Bob",
        "contact": {
            "email": "bob@example.com",
            "phone": "123-456-7890"
        }
    }
}`
var data map[string]any
if err := parser.Decode(json, &data); err != nil {
    fmt.Println("Error:", err)
    return
}

user := data["user"].(map[string]any)
contact := user["contact"].(map[string]any)
fmt.Println(contact["email"]) // Output: bob@example.com
```

#### Array of Objects
```go
json := `{
    "users": [
        {"name": "Alice", "age": 30},
        {"name": "Bob", "age": 25}
    ]
}`
var data map[string]any
if err := parser.Decode(json, &data); err != nil {
    fmt.Println("Error:", err)
    return
}

users := data["users"].([]any)
firstUser := users[0].(map[string]any)
fmt.Println(firstUser["name"]) // Output: Alice
```

#### Nested Arrays
```go
json := `[[1, 2], [3, 4], [5, 6]]`
var matrix []any
if err := parser.Decode(json, &matrix); err != nil {
    fmt.Println("Error:", err)
    return
}

firstRow := matrix[0].([]any)
fmt.Println(firstRow[0]) // Output: 1
```

### Step 6: Complex Real-World Example

```go
json := `{
    "class": 12,
    "section": "A",
    "students": [
        {
            "name": "Alice",
            "marks": {"math": 95, "physics": 88, "chemistry": 92},
            "attendance": 0.95
        },
        {
            "name": "Bob",
            "marks": {"math": 78, "physics": 82, "chemistry": 80},
            "attendance": 0.88
        }
    ],
    "teacher": {
        "name": "Dr. Smith",
        "subjects": ["math", "physics"]
    }
}`

var classData map[string]any
if err := parser.Decode(json, &classData); err != nil {
    fmt.Println("Error:", err)
    return
}

// Access class information
fmt.Println("Class:", classData["class"])
fmt.Println("Section:", classData["section"])

// Access student data
students := classData["students"].([]any)
firstStudent := students[0].(map[string]any)
marks := firstStudent["marks"].(map[string]any)
fmt.Println("First student:", firstStudent["name"])
fmt.Println("Math marks:", marks["math"])

// Access teacher data
teacher := classData["teacher"].(map[string]any)
subjects := teacher["subjects"].([]any)
fmt.Println("Teacher:", teacher["name"])
fmt.Println("First subject:", subjects[0])
```

### Step 7: Error Handling

#### Invalid JSON
```go
json := `{"name": "Alice"` // Missing closing brace
var obj map[string]any
if err := parser.Decode(json, &obj); err != nil {
    fmt.Println("Parse error:", err)
    // Handle the error appropriately
}
```

#### Type Mismatch
```go
json := `"not a number"`
var num float64
if err := parser.Decode(json, &num); err != nil {
    fmt.Println("Type error:", err)
    // Error: cannot assign string to float64
}
```

#### Non-Pointer Argument
```go
json := `{"key": "value"}`
var obj map[string]any
// This will error - must pass a pointer
if err := parser.Decode(json, obj); err != nil {
    fmt.Println("Error:", err)
    // Error: non-nil pointer required
}
```

### Tips

1. Always pass a pointer to `Decode()`: `parser.Decode(json, &variable)`
2. Use `map[string]any` for JSON objects
3. Use `[]any` for JSON arrays
4. Use type assertions to access nested values: `obj["key"].(map[string]any)`
5. Check for errors after every `Decode()` call
6. For null values, use pointer types or interfaces

<details>
  <summary>Automata to identify numbers</summary>
  <img width="564" height="336" alt="image" src="https://github.com/user-attachments/assets/f1c416d9-0b57-41d8-bba7-7e97fe16886f" />
  <br>To identify RFC8259 compliant numbers. It simulates the automata given above. Every character that doesn't have a transition leads to 
  a dead state . Since the language contains every UTF-8 character , it is not possible to show everything . 
  <br>
  <h3>AI Disclaimer</h3>
  I used NanoBanana to convert the hand‑drawn automaton into a digital illustration.
</details>
