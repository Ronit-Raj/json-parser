package scanner

import (
	"unicode"
	"unicode/utf8"
	"fmt"
	"math"
)

type _tokenType uint8

const IDENTIFIER _tokenType = 0
const STRING _tokenType = 1
const NUMBER _tokenType = 2
const BEGIN_ARRAY _tokenType = 3
const BEGIN_OBJECT _tokenType = 4
const END_ARRAY _tokenType = 5
const END_OBJECT _tokenType = 6
const NAME_SEPARATOR _tokenType = 7
const VALUE_SEPARATOR _tokenType = 8
const EOF _tokenType = 9

type Token struct {
	numVal      float64
	stringVal   string
	typeOfToken _tokenType
}
type Error struct {
	msg  string
	code int
}

var Text string
var pointer int

func skipWhiteSpaces() {
	currChar, size := utf8.DecodeRuneInString(Text[pointer:])
	for currChar==' ' || currChar=='\n' || currChar=='\t'{
		currChar,size = utf8.DecodeRuneInString(Text[pointer:])
		pointer+=size
	}
}

func readNumber() (float64,Error) {
	var int float64 = 0
	var exp float64 = 1
	var frac float64 = 0
	var fracMultiplier float64 = 0.1
	var base float64 = 0
	var negBase bool = false
	var negExp bool = false
	var fracBase bool = false
	var state int8 = 0
	var err Error
	for pointer < len(Text) {
		currentChar , charSize := utf8.DecodeRuneInString(Text[pointer:])
		err.msg = fmt.Sprintf("Unexpected chracter %c",currentChar)
		switch state{
		case 0:
			if(currentChar=='0'){
				state = 1
			}else if('1'<=currentChar && currentChar<='9'){
				state = 3
				int*=10
				int+=float64(currentChar-'0')
			}else if(currentChar=='-'){
				negBase = true
				state = 2
			}else{
				state = -1
			}
		case 1:
			if(currentChar=='.'){
				state = 3
				fracBase = true
			}else{
				state = -1
			}
		case 2:
			if(currentChar=='0'){
				state = 1
			}else if('1' <= currentChar && currentChar <='9'){
				state = 3
				int*=10
				int+=float64(currentChar-'0')
			}else{
				state = -1
			}
		case 3:
			if(unicode.IsDigit(currentChar) && !fracBase){
				state = 3
				int*=10
				int+=float64(currentChar-'0')
			}else if(unicode.IsDigit(currentChar) && fracBase){
				state = 3
				frac = frac + (fracMultiplier*(float64(currentChar-'0')))
				fracMultiplier/=10
			}else if(currentChar=='.'){
				state = 4
				fracBase = true
			}else if(currentChar=='e' || currentChar=='E'){
				state = 6
			}else{
				state = -1
			}
		case 4:
			if(unicode.IsDigit(currentChar)){
				state = 5
				frac = frac + (fracMultiplier*(float64(currentChar-'0')))
				fracMultiplier/=10
			}else{
				state = -1
			}
		case 5:
			if(unicode.IsDigit(currentChar)){
				state = 5
				frac = frac + (fracMultiplier*(float64(currentChar-'0')))
				fracMultiplier/=10
			}else if(currentChar=='e' || currentChar=='E'){
				state = 6
			}else {
				state = -1
			}
		case 6:
			if(unicode.IsDigit(currentChar)){
				state = 8
				exp*=10
				exp+=(float64(currentChar-'0'))
			}else if(currentChar=='-'){
				negExp = true
				state = 7
			}else if(currentChar=='+'){
				state = 7
			}else{
				state = -1
			}
		case 7:
			if(unicode.IsDigit(currentChar)){
				state = 8
				exp*=10
				exp+=(float64(currentChar-'0'))
			}else{
				state = -1
			}
		case 8:
			if(unicode.IsDigit(currentChar)){
				state = 8
				exp*=10
				exp+=(float64(currentChar-'0'))
			}else if(currentChar=='.'){
				state = 9
			}else{
				state = -1
			}
		}
		if(state==-1){
			break
		}
		pointer+=charSize
	}

	if(state!=1 || state!=3 || state!=5 || state!=8){
		err.code = -1
	}else{
		err.msg = ""
		err.code = 0
	}

	if(negBase){
		base = -(int+frac)
	}else{
		base = int+frac
	}

	if(negExp){
		exp = 0-exp
	}
	return math.Pow(base,exp) , err 
}

func readString() (string,Error){
	var stringVal string
	startMarker := pointer + 1
	peekPointer := pointer + 1
	err := Error{"Error: expected \" ", -1}
	for peekPointer < len(Text) { //advancing peek pointer to find matching double quotes
		peekChar, pSize := utf8.DecodeLastRuneInString(Text[pointer:])
		if peekChar == '"' && Text[peekPointer-1] != 0x5C {
			/*
				this is the end of a string because we have found a closing double quotes and
				no escape character
			*/
			stringVal = Text[startMarker:peekPointer]
			err = Error{"", 0}
		}
		peekPointer += pSize
	}
	pointer = peekPointer
	return stringVal,err
}

func NextToken() (Token, Error) {
	var currToken Token
	var err Error
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
			currToken = Token{math.NaN(), "", END_OBJECT}
		case rune('}'):
			pointer += size
			currToken = Token{math.NaN(), "", END_ARRAY}
		case rune('"'):
			var stringVal string
			stringVal , err = readString()
			currToken = Token{math.NaN(),stringVal,STRING}		
		case ' ', '\t', '\n', '\r':
			skipWhiteSpaces()
		default:
			if(unicode.IsNumber(currChar) || currChar=='-'){
				var numVal float64
				numVal,err = readNumber()
				currToken = Token{numVal,"",NUMBER}
			}else{
				errorMsg := fmt.Sprintf("Invalid Chracter:%c",currChar)
				err = Error{errorMsg,-1}
			}
		}
		return currToken, err
	}
	return Token{0.0, "", EOF}, Error{"", 0}
}
