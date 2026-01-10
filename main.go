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
			   "array":["hello","world"]}`
	// json := `{}}}`
	var i map[string]any
	if err:=parser.Decode(json,&i); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(i["Name"].(string))
	
}
			