package frontend

import (
	"fmt"
)

type Condition struct {
	Field1   *Token
	Operator *Token
	Field2   *Token
}
type SelectStatement struct {
	Table      *Token
	Columns    []*Token
	Conditions []*Condition
}

var (
	SelectKeyWordToken   = &Token{TokenType: Keyword, Value: "SELECT"}
	FromKeywordToken     = &Token{TokenType: Keyword, Value: "FROM"}
	WhereKeywordToken    = &Token{TokenType: Keyword, Value: "WHERE"}
	AndKeywordToken      = &Token{TokenType: Keyword, Value: "AND"}
	SemicolonSymbolToken = &Token{TokenType: Symbol, Value: ";"}
)

func (t *Token) isEqual(t2 *Token) bool {

	if t.TokenType == t2.TokenType && t.Value == t2.Value {
		return true
	}
	return false
}

var (
// logger = log.New(os.Stdout, "DB4Fun >> ", 0)
)

func ParseSelectStatment(tokenArray []*Token) (statement *SelectStatement, err error) {

	var conditionArray []*Condition
	var columnArray []*Token

	if !SelectKeyWordToken.isEqual(tokenArray[0]) {
		return nil, fmt.Errorf("error at index %d : missing SELECT keyword", tokenArray[0].Position)
	}

	index := 1
	//logger.Printf("started parsing columns at index %d", index)
	index, columnArray, err = ParseColumns(index, tokenArray)
	//logger.Printf("finished parsing columns at index %d", index)
	if err != nil {
		return nil, err
	}
	if index == len(tokenArray) || !FromKeywordToken.isEqual(tokenArray[index]) {
		return nil, fmt.Errorf("error : invalid select statement, missing FROM keyword / comma symbol")
	}
	index++

	if index == len(tokenArray) || tokenArray[index].TokenType != Identifier {
		return nil, fmt.Errorf("error : missing or invalid table name")
	}

	table := tokenArray[index]

	//log.Println("table name : " + table.Value)
	index++

	if index == len(tokenArray) {
		return nil, fmt.Errorf("error : missing semicolon")
	}

	if SemicolonSymbolToken.isEqual(tokenArray[index]) {
		return &SelectStatement{Table: table, Columns: columnArray}, nil
	}

	//log.Println(tokenArray[index].Value)
	if !WhereKeywordToken.isEqual(tokenArray[index]) {
		return nil, fmt.Errorf("error at index %d : missing WHERE keyword", tokenArray[index].Position)
	}

	index++
	//logger.Printf("started parsing conditions at index %d", index)
	index, conditionArray, err = ParseConditions(index, tokenArray)
	//logger.Printf("finished parsing conditions at index %d", index)
	if err != nil {
		return nil, err
	}

	index++

	return &SelectStatement{Table: table, Columns: columnArray, Conditions: conditionArray}, nil
}

func ParseConditions(index int, tokenArray []*Token) (newIndex int, conditionArray []*Condition, err error) {

	conditionArray = make([]*Condition, 0)

	currIndex := index
	//logger.Printf("condition parsing begins at index %d", currIndex)
	for currIndex < len(tokenArray) {
		//logger.Printf("new condition being parsed, starting from index %d", currIndex)

		newCondition := &Condition{}
		for i := 0; i < 3; i++ {

			if currIndex+i >= len(tokenArray) {
				//logger.Printf("index %d out of bound of token array with length %d", (currIndex + i), len(tokenArray))
				return index, nil, fmt.Errorf("error : invalid select statement, missing / invalid condition after WHERE clause")
			}
			currToken := tokenArray[currIndex+i]
			if i == 0 {
				if currToken.TokenType != Identifier {
					//logger.Printf("index %d out of bound of token array with length %d", (currIndex + i), len(tokenArray))
					return index, nil, fmt.Errorf("error at index %d : invalid column name", currToken.Position)
				}
				newCondition.Field1 = currToken
			}
			if i == 1 {
				if currToken.TokenType != Operator {
					return index, nil, fmt.Errorf("error at index %d : invalid operator", currToken.Position)
				}
				newCondition.Operator = currToken

			}
			if i == 2 {
				if currToken.TokenType != String {
					return index, nil, fmt.Errorf("error at index %d : invalid constant value", currToken.Position)
				}
				newCondition.Field2 = currToken
			}

		}
		currIndex += 3
		conditionArray = append(conditionArray, newCondition)
		if currIndex == len(tokenArray) {
			return index, nil, fmt.Errorf("error : missing semicolon")
		}

		if SemicolonSymbolToken.isEqual(tokenArray[currIndex]) {
			//logger.Printf("finished parsing conditions at index %d", currIndex)
			return currIndex, conditionArray, nil
		}
		if !AndKeywordToken.isEqual(tokenArray[currIndex]) {
			return index, nil, fmt.Errorf("error at index %d : missing AND keyword between conditions", tokenArray[currIndex].Position)
		}
		currIndex++
	}

	return currIndex, conditionArray, nil

}

func ParseColumns(index int, tokenArray []*Token) (newIndex int, columnArray []*Token, err error) {

	columnArray = make([]*Token, 0)
	endOfColumnList := false

	currIndex := index
	for currIndex < len(tokenArray) && !endOfColumnList {

		currToken := tokenArray[currIndex]

		if currIndex%2 == 1 {
			if currToken.TokenType == Symbol && currToken.Value == "*" {
				endOfColumnList = true

			} else if currToken.TokenType != Identifier {

				return index, nil, fmt.Errorf("error at index %d : invalid column name", currToken.Position)

			}

			columnArray = append(columnArray, currToken)
			currIndex++

		} else {

			if currToken.TokenType != Symbol && currToken.Value != "," {
				endOfColumnList = true
			} else {
				currIndex++
			}
		}

	}

	return currIndex, columnArray, nil
}
