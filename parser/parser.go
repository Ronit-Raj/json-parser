package parser

import (
	"fmt"
	"json_parser/scanner"
)



func Decode(json string) (any,error) {
	scanner.Text = json
	err,val := value()
	if err != nil {
		return nil,err
	}
	return val,nil
}

func value() (error,any) {
	for token,err := scanner.NextToken() ; token.TypeOfToken != scanner.EOF ; token,err = scanner.NextToken() {
		if err != nil {
			return err , nil 
		}

		switch token.TypeOfToken {
		case scanner.LITERAL_TRUE:
			return nil,true 
		case scanner.LITERAL_FALSE:
			return nil,false
		case scanner.LITERAL_NULL:
			return nil,scanner.LITERAL_NULL
		case scanner.STRING:
			return nil,token.StringVal 
		case scanner.NUMBER:
			return nil,token.NumVal 
		case scanner.BEGIN_ARRAY:
			err,obj := array()
			if err != nil{
				return err,nil
			}
			return nil,obj
		case scanner.BEGIN_OBJECT:
			err,arr := member()
			if err != nil {
				return err,nil
			}
			return nil,arr
		}
	}
	return nil,nil
}
func member() (error,map[string]any){
	decodedObj := make(map[string]any)
	type state int8
	const (
		start state = iota
		end
		parsedKey
		parsedNameSep
		parsedValue
		parsedValSep
	)
	var st state 
	st = start
	var currentKey string
	for token,err := scanner.NextToken() ; token.TypeOfToken != scanner.EOF ; token,err = scanner.NextToken(){
		if err != nil {
			return err , nil
		}

		switch st {
		case start:
			if token.TypeOfToken == scanner.END_OBJECT{
				st = end
				return nil,decodedObj
			}else if token.TypeOfToken == scanner.STRING { 
				currentKey = token.StringVal
				st = parsedKey
			}else{
				return fmt.Errorf(`Error:Expected string or "}" inside object`) , nil
			}
		case parsedKey:
			if token.TypeOfToken == scanner.NAME_SEPARATOR{
				// st = parsedNameSep
				err , val := value()
				if err!=nil {
					return err,nil
				}
				decodedObj[currentKey] = val
				st = parsedValue
			}else{
				return fmt.Errorf(`Error:Expected ":" after string `) , nil
			}
		// case parsedNameSep:
		// 	err , val := value()
		// 	if err!=nil {
		// 		return err,nil
		// 	}
		// 	decodedObj[currentKey] = val
		// 	st = parsedValue
		case parsedValue:
			if token.TypeOfToken == scanner.END_OBJECT{
				st = end 
				return nil,decodedObj
			}else if token.TypeOfToken == scanner.VALUE_SEPARATOR{
				st = parsedValSep
			}else{
				return fmt.Errorf(`Error:Unexpected end of object`) , nil
			}
		case parsedValSep:
			if token.TypeOfToken == scanner.STRING {
				st = parsedKey
				currentKey = token.StringVal
			}else{
				return fmt.Errorf(`Error:Expected string `),nil
			}
		}
		
	}
	return nil , decodedObj
}

func array() (error,[]any){
	decodedArr := make([]any,1)
	type state int8
	const (
		start state = iota 
		end 
		parsedVal 
		parsedValSep
	)
	var st state
	st = start
	for token,err := scanner.NextToken() ; token.TypeOfToken != scanner.EOF ; token,err = scanner.NextToken(){
		if err != nil {
			return err , nil
		}

		switch st {
		case start:
			err,val := value()
			if token.TypeOfToken == scanner.END_ARRAY {
				st = end
				return err,decodedArr
			}else{
				if(err!=nil){
					return err,nil
				}
				decodedArr = append(decodedArr, val)
				st = parsedVal
			}
		case parsedVal:
			if token.TypeOfToken == scanner.END_ARRAY {
				st = end
				return err,decodedArr
			}else if token.TypeOfToken == scanner.VALUE_SEPARATOR {
				st = parsedValSep
			}
		case parsedValSep:
			err,val := value()
			if err != nil {
				return err,nil
			}
			decodedArr = append(decodedArr, val)
			st = parsedVal
		}
	}
	return nil,decodedArr
}