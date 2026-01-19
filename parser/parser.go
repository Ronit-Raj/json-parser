package parser

import (
	"fmt"
	"json_parser/scanner"
	"reflect"
)

func Decode(json string, decodeVal any) error {
	rv := reflect.ValueOf(decodeVal)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("Decode expects a pointer")
	}
	scanner.Text = json
	err, val := value()
	if err != nil {
		return err
	}

	dst := rv.Elem()
	src := reflect.ValueOf(val)
	if !src.Type().AssignableTo(dst.Type()) {
		return fmt.Errorf("cannot assign %v to %v", src.Type(), dst.Type())
	}
	dst.Set(src)
	return nil
}

func value() (error, any) {
	for token, err := scanner.PeekToken(); token.TypeOfToken != scanner.EOF; token, err = scanner.PeekToken() {
		if err != nil {
			return err, nil
		}

		switch token.TypeOfToken {
		case scanner.NUMBER:
			scanner.NextToken() // consume the token
			return nil, token.NumVal
		case scanner.STRING:
			scanner.NextToken()
			return nil, token.StringVal
		case scanner.LITERAL_FALSE:
			scanner.NextToken()
			return nil, false
		case scanner.LITERAL_NULL:
			scanner.NextToken()
			return nil, nil
		case scanner.LITERAL_TRUE:
			scanner.NextToken()
			return nil, true
		case scanner.BEGIN_ARRAY:
			return array()
		case scanner.BEGIN_OBJECT:
			return member()
		default:
			return fmt.Errorf(`Error: Unexpected token `), nil
		}
	}
	return nil, nil // this should be unreachable
}
func member() (error, map[string]any) {
	scanner.NextToken() // consume '{'
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
	for token, err := scanner.NextToken(); token.TypeOfToken != scanner.EOF; token, err = scanner.NextToken() {
		if err != nil {
			return err, nil
		}

		switch st {
		case start:
			if token.TypeOfToken == scanner.END_OBJECT {
				st = end
				return nil, decodedObj
			} else if token.TypeOfToken == scanner.STRING {
				currentKey = token.StringVal
				st = parsedKey
			} else {
				return fmt.Errorf(`Error:Expected string or "}" inside object`), nil
			}
		case parsedKey:
			if token.TypeOfToken == scanner.NAME_SEPARATOR {
				err, val := value()
				if err != nil {
					return err, nil
				}
				decodedObj[currentKey] = val
				st = parsedValue
			} else {
				return fmt.Errorf(`Error:Expected ":" after string `), nil
			}
		case parsedValue:
			if token.TypeOfToken == scanner.END_OBJECT {
				st = end
				return nil, decodedObj
			} else if token.TypeOfToken == scanner.VALUE_SEPARATOR {
				st = parsedValSep
			} else {
				return fmt.Errorf(`Error:Unexpected end of object`), nil
			}
		case parsedValSep:
			if token.TypeOfToken == scanner.STRING {
				st = parsedKey
				currentKey = token.StringVal
			} else {
				return fmt.Errorf(`Error:Expected string `), nil
			}
		}

	}
	return nil, decodedObj
}

func array() (error, []any) {
	scanner.NextToken() // consume '['
	decodedArr := make([]any, 0)
	type state int8
	const (
		start state = iota
		end
		parsedVal
		parsedValSep
	)
	var st state
	st = start
	for token, err := scanner.PeekToken(); token.TypeOfToken != scanner.EOF; token, err = scanner.PeekToken() {
		if err != nil {
			return err, nil
		}

		switch st {
		case start:
			err, val := value()
			if token.TypeOfToken == scanner.END_ARRAY {
				scanner.NextToken()
				st = end
				return err, decodedArr
			} else {
				if err != nil {
					return err, nil
				}
				decodedArr = append(decodedArr, val)
				st = parsedVal
			}
		case parsedVal:
			if token.TypeOfToken == scanner.END_ARRAY {
				scanner.NextToken()
				st = end
				return err, decodedArr
			} else if token.TypeOfToken == scanner.VALUE_SEPARATOR {
				scanner.NextToken()
				st = parsedValSep
			}
		case parsedValSep:
			err, val := value()
			if err != nil {
				return err, nil
			}
			decodedArr = append(decodedArr, val)
			st = parsedVal
		}
	}
	return nil, decodedArr
}
