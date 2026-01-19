package scanner

import (
	"unicode"
	"unicode/utf8"
	"fmt"
	"math"
	"strconv"
)

type TokenType uint8

const STRING TokenType = 1
const NUMBER TokenType = 2
const BEGIN_ARRAY TokenType = 3
const BEGIN_OBJECT TokenType = 4
const END_ARRAY TokenType = 5
const END_OBJECT TokenType = 6
const NAME_SEPARATOR TokenType = 7
const VALUE_SEPARATOR TokenType = 8
const LITERAL_FALSE TokenType = 9
const LITERAL_TRUE TokenType = 10
const LITERAL_NULL TokenType = 11
const EOF TokenType = 12

type Token struct {
	NumVal      float64
	StringVal   string
	TypeOfToken TokenType
}
type SyntaxError struct {
	Msg  string
	Position int
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("Error:%d %s",e.Position,e.Msg)
}

var Text string
var pointer int

func ResetPointer() {
	pointer = 0
}

func skipWhiteSpaces() {
	currChar, size := utf8.DecodeRuneInString(Text[pointer:])
	for currChar==' ' || currChar=='\n' || currChar=='\t'{
		pointer+=size
		currChar,size = utf8.DecodeRuneInString(Text[pointer:])
	}
}

// check https://github.com/Ronit-Raj/json-parser/blob/main/README.md for automata
func readNumber() (float64,error) {
	var state int8 = 0
	var start int = pointer
	var err SyntaxError
	loop:
	for pointer < len(Text) {
		currentChar , charSize := utf8.DecodeRuneInString(Text[pointer:])
		err.Msg = fmt.Sprintf("Unexpected chracter %c",currentChar)
		switch state{
		case 0:
			if(currentChar=='0'){
				state = 1
			}else if('1'<=currentChar && currentChar<='9'){
				state = 3
			}else if(currentChar=='-'){
				state = 2
			}else{
				state = -1
			}
		case 1:
			if(currentChar=='.'){
				state = 3
			}else if(currentChar=='e' || currentChar=='E'){
				state = 6
			}else{
				break loop
			}
		case 2:
			if(currentChar=='0'){
				state = 1
			}else if('1' <= currentChar && currentChar <='9'){
				state = 3
			}else{
				state = -1
			}
		case 3:
			if(unicode.IsDigit(currentChar)){
				state = 3
			}else if(unicode.IsDigit(currentChar)){
				state = 3
			}else if(currentChar=='.'){
				state = 4
			}else if(currentChar=='e' || currentChar=='E'){
				state = 6
			}else{
				break loop
			}
		case 4:
			if(unicode.IsDigit(currentChar)){
				state = 5
			}else{
				state = -1
			}
		case 5:
			if(unicode.IsDigit(currentChar)){
				state = 5
			}else if(currentChar=='e' || currentChar=='E'){
				state = 6
			}else {
				break loop
			}
		case 6:
			if(unicode.IsDigit(currentChar)){
				state = 8
			}else if(currentChar=='-'){
				state = 7
			}else if(currentChar=='+'){
				state = 7
			}else{
				state = -1
			}
		case 7:
			if(unicode.IsDigit(currentChar)){
				state = 8
			}else{
				state = -1
			}
		case 8:
			if(unicode.IsDigit(currentChar)){
				state = 8
			}else if(currentChar=='.'){
				state = 9
			}else{
				break loop
			}
		}
		if(state==-1){
			break
		}
		pointer+=charSize
	}

	if(state!=1 && state!=3 && state!=5 && state!=8){
		err.Position = pointer
		return math.NaN(),err
	}else{
		num,_ := strconv.ParseFloat(Text[start:pointer],64)
		return num , nil
	}
}

func readString() (string,error){
	var stringVal string
	startMarker := pointer 
	peekPointer := pointer 
	for peekPointer < len(Text) { //advancing peek pointer to find matching double quotes
		peekChar, pSize := utf8.DecodeRuneInString(Text[peekPointer:])
		if peekChar == '"' && Text[peekPointer-1] != 0x5C {
			/*
				this is the end of a string because we have found a closing double quotes and
				no escape character
			*/
			stringVal = Text[startMarker:peekPointer]
			peekPointer += pSize
			pointer = peekPointer
			return stringVal , nil
		}
		peekPointer += pSize
	}
	pointer = peekPointer
	return "",SyntaxError{
		Msg: "unterminated string",
		Position: startMarker,
	}
}

func match(lex string) (bool){
	for _,val := range lex {
		if rune(Text[pointer]) != val {
			return false
		}
		pointer++
	}
	return true
}

func PeekToken() (Token,error){
	peekPointer := pointer
	peekToken,err := NextToken()
	if(err!=nil){
		pointer = peekPointer
		return Token{0.0, "", EOF},err
	}
	pointer = peekPointer
	return peekToken,nil
}

func NextToken() (Token, error) {
	var currToken Token
	var err error
	if pointer < len(Text) {
		currChar, size := utf8.DecodeRuneInString(Text[pointer:])

		switch currChar {
		case rune(':'):
			pointer += size
			currToken = Token{math.NaN(), "", NAME_SEPARATOR}
		case rune(','):
			pointer += size
			currToken = Token{math.NaN(), "", VALUE_SEPARATOR}
		case rune('{'):
			pointer += size
			currToken = Token{math.NaN(), "", BEGIN_OBJECT}
		case rune('['):
			pointer += size
			currToken = Token{math.NaN(), "", BEGIN_ARRAY}
		case rune(']'):
			pointer += size
			currToken = Token{math.NaN(), "", END_ARRAY}
		case rune('}'):
			pointer += size
			currToken = Token{math.NaN(), "", END_OBJECT}
		case rune('"'):
			var stringVal string
			pointer += size
			stringVal , err = readString()
			currToken = Token{math.NaN(),stringVal,STRING}
		case rune('f'):
			if (match("false")){
				currToken = Token{math.NaN(),"",LITERAL_FALSE}
			}else {
				errorMsg := fmt.Sprintf("Invalid Chracter:%c",currChar)
				err = SyntaxError{errorMsg,pointer}
			}
		case rune('t'):
			if (match("true")){
				currToken = Token{math.NaN(),"",LITERAL_TRUE}
			}else{
				errorMsg := fmt.Sprintf("Invalid Chracter:%c",currChar)
				err = SyntaxError{errorMsg,pointer}
			}
		case rune('n'):
			if (match("null")){
				currToken = Token{math.NaN(),"",LITERAL_NULL}
			}else{
				errorMsg := fmt.Sprintf("Invalid Chracter:%c",currChar)
				err = SyntaxError{errorMsg,pointer}
			}
		case ' ', '\t', '\n', '\r':
			skipWhiteSpaces()
			currToken,err = NextToken()
		default:
			if(unicode.IsNumber(currChar) || currChar=='-'){
				var numVal float64
				numVal,err = readNumber()
				currToken = Token{numVal,"",NUMBER}
			}else{
				errorMsg := fmt.Sprintf("Invalid Chracter:%c",currChar)
				err = SyntaxError{errorMsg,pointer}
			}
		}
		return currToken, err
	}
	return Token{0.0, "", EOF}, nil
}
