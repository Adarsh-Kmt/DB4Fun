package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Adarsh-Kmt/DB4Fun/backend/btree"
	"github.com/Adarsh-Kmt/DB4Fun/frontend"
)

var (
	logger = log.New(os.Stdout, "DB4Fun >> ", 0)
)

func main() {

	scanner := bufio.NewReader(os.Stdin)
	fmt.Println("DB4Fun")
	fmt.Println("1) enter 'exit' to exit terminal")
	fmt.Println("2) enter 'test-btree' to test working of b tree with maximum 3 items per node.")
	fmt.Println("3) enter 'test-frontend' to test SQL tokenizer and parser. ")
	for {
		logger.Print("enter input: ")
		input, err := scanner.ReadString('\n')
		input = strings.TrimRight(input, "\r\n")
		logger.Println("string received : " + input)
		if err != nil {
			return
		}
		if input == "exit" {
			logger.Println("exiting...")
			return
		}

		if input == "test-btree" {

			logger.Println("testing b-tree...")
			bTree := btree.BTreeInit(3)
			ok := true
			for ok {
				logger.Print("1) enter 'insert' to insert item     2) enter 'search' to search for item")

				option, err := scanner.ReadString('\n')
				option = strings.TrimRight(option, "\r\n")
				if err != nil {
					continue
				}

				logger.Print("enter key : ")
				key, err := scanner.ReadString('\n')
				key = strings.TrimRight(key, "\r\n")
				if err != nil {
					return
				}
				if option == "insert" {

					logger.Print("enter value : ")
					value, err := scanner.ReadString('\n')
					value = strings.TrimRight(value, "\r\n")
					if err != nil {
						return
					}
					bTree.InsertItem(key, value)

				} else {

					value, found := bTree.SearchForItem(key)
					if found {
						logger.Printf("value for key %s = %s", key, value)
					} else {
						logger.Printf("key %s not found in btree.", key)
					}
				}

				logger.Print("continue? (Y/N)")

				YNResponse, err := scanner.ReadString('\n')
				YNResponse = strings.TrimRight(YNResponse, "\r\n")
				if err != nil {
					return
				}

				if strings.Compare(YNResponse, "N") == 0 {
					ok = false
				}
			}
		}
		if input == "test-frontend" {

			logger.Println("enter SQL statement : ")
			str, err := scanner.ReadString('\n')
			str = strings.TrimRight(str, "\r\n")
			if err != nil {
				return
			}
			tokenArray, err := frontend.Tokenizer(str)

			if err != nil {
				logger.Printf("%s", err.Error())

			}

			DisplayTokenizerOutput(tokenArray)

			selectStatement, err := frontend.ParseSelectStatment(tokenArray)

			if err != nil {
				logger.Println(err.Error())

			}

			DisplayParserOutput(selectStatement)

		}

	}

	// if true {

	// 	bTree := btree.BTreeInit(4)
	// 	bTree.InsertItem("4", "a")
	// 	bTree.InsertItem("5", "b")
	// 	bTree.InsertItem("6", "c")
	// 	bTree.InsertItem("9", "d")
	// 	bTree.InsertItem("10", "e")
	// 	bTree.InsertItem("3", "f")
	// 	bTree.InsertItem("2", "g")
	// 	bTree.InsertItem("8", "h")
	// 	bTree.InsertItem("7", "i")
	// 	bTree.InsertItem("11", "j")
	// 	bTree.InsertItem("1", "k")
	// 	bTree.InsertItem("0", "l")
	// 	//backend.BTreeTraversal(bTree.Root)
	// 	return
	// }
	// str := "SELECT first_name, last_name FROM user_table WHERE first_name = 'adarsh' AND last_name = 'kamath';"
	// tokenArray, err := frontend.Tokenizer(str)

	// if err != nil {
	// 	logger.Printf("%s", err.Error())
	// 	return
	// }

	// DisplayTokenizerOutput(tokenArray)

	// selectStatement, err := frontend.ParseSelectStatment(tokenArray)

	// if err != nil {
	// 	logger.Println(err.Error())
	// 	return
	// }

	// DisplayParserOutput(selectStatement)

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
