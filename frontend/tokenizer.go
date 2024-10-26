package frontend

import (
	"fmt"
)

type TokenType int

const (
	Keyword TokenType = iota
	Identifier
	Symbol
	String
	Operator
)

type KeywordType string

const (
	SelectKeyword KeywordType = "SELECT"
	CreateKeyword KeywordType = "CREATE"
	InsertKeyword KeywordType = "INSERT"
	FromKeyword   KeywordType = "FROM"
	TableKeyword  KeywordType = "TABLE"
	IntoKeyword   KeywordType = "INTO"
	ValuesKeyword KeywordType = "VALUES"
	WhereKeyword  KeywordType = "WHERE"
)

var (
	keywordArray     = []KeywordType{SelectKeyword, CreateKeyword, IntoKeyword, FromKeyword, TableKeyword, InsertKeyword, ValuesKeyword, WhereKeyword}
	SymbolArray      = []SymbolType{OpeningParenthSymbol, ClosingParenthSymbol, CommaSymbol, SemicolonSymbol, StarSymbol}
	TwoOperatorArray = []OperatorType{GreaterThanEqualOperator, LesserThanEqualOperator, NotEqualOperator}
	OneOperatorArray = []OperatorType{LesserThanOperator, GreaterThanOperator, EqualOperator}
)

type SymbolType string

const (
	OpeningParenthSymbol SymbolType = "("
	ClosingParenthSymbol SymbolType = ")"
	CommaSymbol          SymbolType = ","
	StarSymbol           SymbolType = "*"
	SemicolonSymbol      SymbolType = ";"
)

type OperatorType string

const (
	GreaterThanEqualOperator OperatorType = ">="
	LesserThanEqualOperator  OperatorType = "<="
	LesserThanOperator       OperatorType = "<"
	GreaterThanOperator      OperatorType = ">"
	EqualOperator            OperatorType = "="
	NotEqualOperator         OperatorType = "<>"
)

type Token struct {
	TokenType TokenType
	Value     string
	Position  int
}

type searchFor func(str string, cursor *int) (token *Token, ok bool)

func searchForKeyword(str string, cursor *int) (token *Token, ok bool) {
	//log.Printf("called check for keyword method at index %d", *cursor)

	for _, key := range keywordArray {

		if *cursor+len(key) < len(str) && str[*cursor:*cursor+len(key)] == string(key) {
			//log.Println(string(key))
			newToken := &Token{TokenType: Keyword, Value: string(key), Position: *cursor}
			*cursor += len(key)
			return newToken, true
		}
	}
	//log.Println("not a keyword")
	return nil, false
}

func searchForIdentifier(str string, cursor *int) (token *Token, ok bool) {

	//log.Printf("called check for Identifier method at index %d", *cursor)

	identifier := ""
	tempCursor := *cursor

	for tempCursor < len(str) && str[tempCursor] != ' ' && str[tempCursor] != ',' && str[tempCursor] != '(' && str[tempCursor] != ')' && str[tempCursor] != ';' {

		identifier = identifier + string(str[tempCursor])
		tempCursor++
	}

	if len(identifier) != 0 {
		//log.Println(Identifier)
		newToken := &Token{TokenType: Identifier, Value: identifier, Position: *cursor}
		*cursor = tempCursor
		return newToken, true
	}

	return nil, false
}

func searchForSymbol(str string, cursor *int) (token *Token, ok bool) {
	//log.Printf("called check for symbol method at index %d", *cursor)

	for _, sym := range SymbolArray {

		if str[*cursor:*cursor+1] == string(sym) {
			symbol := string(str[*cursor])
			//log.Println(symbol)
			newToken := &Token{TokenType: Symbol, Value: symbol, Position: *cursor}
			(*cursor)++
			return newToken, true
		}
	}
	//log.Println("not a symbol.")
	return nil, false
}
func searchForString(str string, cursor *int) (token *Token, ok bool) {

	//log.Printf("called check for string method at index %d", *cursor)
	if *cursor < len(str) && str[*cursor:*cursor+1] != "'" {

		return nil, false
	}

	strToken := ""
	tempCursor := *cursor + 1

	for tempCursor < len(str) && str[tempCursor:tempCursor+1] != "'" {
		//log.Printf("called check for string method at index %d", tempCursor)
		strToken = strToken + string(str[tempCursor])
		//log.Println(strToken)
		//time.Sleep(2 * time.Second)
		tempCursor++
	}
	newToken := &Token{TokenType: String, Value: strToken, Position: *cursor}
	*cursor = tempCursor + 1
	//log.Println("exited string method")

	return newToken, true
}

func searchForOperator(str string, cursor *int) (token *Token, ok bool) {

	if *cursor+2 < len(str) {

		for _, op := range TwoOperatorArray {

			if str[*cursor:*cursor+2] == string(op) {
				newToken := &Token{TokenType: Operator, Value: string(op), Position: *cursor}
				(*cursor) += 2
				return newToken, true
			}
		}
	}

	if *cursor+1 < len(str) {

		for _, op := range OneOperatorArray {

			if str[*cursor:*cursor+1] == string(op) {
				newToken := &Token{TokenType: Operator, Value: string(op), Position: *cursor}
				(*cursor) += 1
				return newToken, true
			}
		}
	}

	return nil, false
}
func Tokenizer(str string) ([]*Token, error) {

	arr := make([]*Token, 0)

	cursor := 0

	searchForFuncArray := []searchFor{searchForString, searchForKeyword, searchForSymbol, searchForOperator, searchForIdentifier}

	for cursor != len(str) {

		for cursor < len(str) && str[cursor] == ' ' {
			cursor++
		}

		var ok bool
		ok = false
		for _, search := range searchForFuncArray {

			var nextToken *Token
			if nextToken, ok = search(str, &cursor); ok {
				//log.Println("token formed with value : " + nextToken.Value)
				arr = append(arr, nextToken)
				break
			}
		}

		if !ok {
			return nil, fmt.Errorf("error at character %d", cursor)
		}

	}

	return arr, nil
}
