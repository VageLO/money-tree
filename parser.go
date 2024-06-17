package main

import (
	"fmt"
	"regexp"

	"github.com/dslipak/pdf"
	"github.com/fatih/color"
)

var (
	filename       = "./alfa.pdf"
	Price          = regexp.MustCompile(`(\d+\.\d{2} [A-Z]{3})`)
	DateTime       = regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
	TransactionNum = regexp.MustCompile(`\d{16}`)
)

type Transaction struct {
	id          string
	date        string
	typeof      string
	status      string
	price       string
	discription string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	r, err := pdf.Open(filename)
	check(err)

	var transactions []Transaction

	for pageNum := 1; pageNum <= r.NumPage(); pageNum++ {
		fmt.Printf("Page: %d\n", pageNum)

		page := r.Page(pageNum)

		var sentence pdf.Text
		var temp pdf.Text

		var transaction Transaction

		//var tempID, tranID, priceID int

		texts := page.Content().Text
		for _, text := range texts {
			if transaction.id == "" {
				fmt.Println("CALL")
				if text.Y == sentence.Y {
					sentence.S = sentence.S + text.S
				} else {
					sentence = text
				}
			} else {
				temp = text
				temp.S = temp.S + text.S
			}
			//} else if sentence.X == text.X {
			//	temp = text
			//	temp.S = sentence.S + " " + temp.S
			//	sentence = temp
			//} else {
			//	sentence = text
			//}

			fmt.Printf(color.GreenString("TEXT: %v "), temp.S)
			fmt.Printf(color.RedString("SENTENCE: %v\n"), sentence)
			switch {
			//case Price.MatchString(sentence.S):
			//	priceID = i + 1
			//	//fmt.Println(tempID, tranID, i, Price.FindString(sentence.S))
			//	Parse(texts[tempID:i-(len(sentence.S)-1)], &transaction)

			//	transaction.price = Price.FindString(sentence.S)
			//	transactions = append(transactions, transaction)
			//	sentence.S = ""
			//case DateTime.MatchString(sentence.S):
			//	tempID = i + 1
			//	//fmt.Println(DateTime.FindString(sentence.S))
			//	transaction.date = DateTime.FindString(sentence.S)
			//	sentence.S = ""
			case TransactionNum.MatchString(sentence.S):
				//	tranID = i
				fmt.Println(TransactionNum.FindString(sentence.S))
				transaction.id = TransactionNum.FindString(sentence.S)
				transactions = append(transactions, transaction)
				transaction.id = ""
				//	if priceID != 0 && i+1 > priceID {
				//		fmt.Println(texts[priceID : tranID-(len(sentence.S)-1)])
			}
			//	sentence.S = ""

			//}
		}
	}
	fmt.Println(transactions)
}

func Parse(array []pdf.Text, transaction *Transaction) {
	var sentence pdf.Text
	var temp pdf.Text

	for i, text := range array {
		if text.Y == sentence.Y {
			sentence.S = sentence.S + text.S
			if len(array)-1 == i {
				transaction.status = sentence.S
			}
		} else if sentence.X == text.X {
			temp = text
			temp.S = sentence.S + " " + temp.S
			sentence = temp
		} else {
			transaction.typeof = sentence.S
			sentence = text
		}
	}
}
