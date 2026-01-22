package main

import (
	"fmt"
	"json_parser/parser"
)

func main() {
	json := `{"class":12,"sec":"A","Name":"ronit",
			   "marks":{"phy":90,"chem":85,"maths":90},
			   "co-cirrcular":{},
			   "address":null,
			   "array":["hello","world"],
			   "empty-array":[]}`
	// json := `{}}}`
	var i map[string]any
	if err := parser.Decode(json, &i); err != nil {
		fmt.Println(err)
		return
	}
	
	marks, ok := i["marks"].(map[string]any)
	if !ok {
		fmt.Println("marks field is not a map")
		return
	}
	
	chemMarks, ok := marks["chem"]
	if !ok {
		fmt.Println("chem field not found")
		return
	}
	
	fmt.Println(chemMarks)

	json = `["💀12",
	23  ,
	true]`
	var j []any
	if err := parser.Decode(json, &j); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(j[1])

	json = `{"users": [{"name": "Alice", "age": 30}, {"name": "Bob", "age": 25}]}`
	var k map[string]any
	if err := parser.Decode(json, &k); err != nil {
		fmt.Println(err)
		return
	}
	users := k["users"].([]any)
	secondUser := users[1].(map[string]any)
	fmt.Println(secondUser["name"])

}
