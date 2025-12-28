package main 

import "fmt"
import "json-parser/scanner"

func main() {
	scanner.Text = "-3.67e2}"
	token , err := scanner.NextToken()
	for token.TypeOfToken != scanner.EOF {
		fmt.Printf("string-val = %s num-val =%f type=%d \n",token.StringVal,token.NumVal,token.TypeOfToken)
		fmt.Printf("error-message=%s error-code=%d \n",err.Msg,err.Code)
		token , err = scanner.NextToken()
	}
}