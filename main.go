package main

import (
	"log"
	"os"

	"github.com/Adarsh-Kmt/DB4Fun/frontend"
)

func main() {

	str := "SELECT first_name, last_name FROM user_table WHERE first_name = 'adarsh' AND last_name = 'kamath';"
	logger := log.New(os.Stdout, "DB4Fun >> ", 0)
	arr, err := frontend.Tokenizer(str)

	if err != nil {
		logger.Printf("%s", err.Error())
		return
	}

	logger.Println()
	logger.Println("-------- TOKENIZER -------")
	logger.Println()

	for ind, token := range arr {

		logger.Printf("[ index : %d token code : %d , value :  %s  position : %d ]\n", ind, token.TokenType, token.Value, token.Position)
	}
	//logger.Println("END")

	selectStatement, err := frontend.ParseSelectStatment(arr)

	if err != nil {
		logger.Println(err.Error())
		return
	}
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
