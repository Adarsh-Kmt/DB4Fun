package main

import (
	"log"
	"os"

	"github.com/Adarsh-Kmt/DB4Fun/frontend"
)

var (
	logger = log.New(os.Stdout, "DB4Fun >> ", 0)
)

func main() {

	str := "SELECT first_name, last_name FROM user_table WHERE first_name = 'adarsh' AND last_name = 'kamath';"
	tokenArray, err := frontend.Tokenizer(str)

	if err != nil {
		logger.Printf("%s", err.Error())
		return
	}

	DisplayTokenizerOutput(tokenArray)

	selectStatement, err := frontend.ParseSelectStatment(tokenArray)

	if err != nil {
		logger.Println(err.Error())
		return
	}

	DisplayParserOutput(selectStatement)

}

func DisplayTokenizerOutput(tokenArray []*frontend.Token) {

	logger.Println()
	logger.Println("-------- TOKENIZER -------")
	logger.Println()

	for ind, token := range tokenArray {

		logger.Printf("[ index : %d token code : %d , value :  %s  position : %d ]\n", ind, token.TokenType, token.Value, token.Position)
	}
}

func DisplayParserOutput(selectStatement *frontend.SelectStatement) {

	logger.Println()
	logger.Println("---- SELECT STATEMENT ----")
	logger.Println()
	logger.Println("table : " + selectStatement.Table.Value)
	logger.Println()
	for ind, column := range selectStatement.Columns {
		logger.Printf("column %d : %s", ind+1, column.Value)
	}
	logger.Println()
	for ind, condition := range selectStatement.Conditions {
		logger.Printf("condition %d : %s %s %s", ind+1, condition.Field1.Value, condition.Operator.Value, condition.Field2.Value)
	}

	logger.Println()
	logger.Println("-------------------------")
}
