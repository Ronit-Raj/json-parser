package main 

import "fmt"
import "unicode/utf8"
import "scanner"

func main() {
	s := "a世界😅😅b 45"
	i := 0 
	fmt.Println(utf8.RuneCountInString(s))
	for i<len(s){
		char , size := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%c",char)
		i+=size
	}
}