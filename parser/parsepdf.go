package parser

import (
	"os"
	"regexp"
	"strings"

	"github.com/dslipak/pdf"
)

var (
	PriceWithAcronym = regexp.MustCompile(`(\d+\.\d{2}\s*[A-Z]{3})`)
	Price            = regexp.MustCompile(`\d+\.\d{2}`)
	Acronym          = regexp.MustCompile(`[A-Z]{3}`)
	Date             = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	Time             = regexp.MustCompile(`\d{2}:\d{2}:\d{2}`)
	TransactionNum   = regexp.MustCompile(`\d{16}`)
	TransactionType  = regexp.MustCompile(`[А-Я][^А-Я]*[А-Я]`)
)

type Transaction struct {
	Id          string
	Date        string
	Time        string
	Typeof      string
	Status      string
	Price       string
	Acronym     string
	Description string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ParsePdf(filename string) []Transaction {
	r, err := pdf.Open(filename)
	check(err)

	var transactions []Transaction

	for pageNum := 1; pageNum <= r.NumPage(); pageNum++ {

		page := r.Page(pageNum)

		var sentence pdf.Text
		var temp pdf.Text

		var transaction Transaction

		var tranID int

		var temp_str string

		texts := page.Content().Text
		for i, text := range texts {
			if text.Y == sentence.Y {
				sentence.S = sentence.S + text.S
			} else if sentence.X == text.X {
				temp = text
				temp.S = sentence.S + " " + temp.S
				sentence = temp
			} else {
				sentence = text
			}

			if i+1 == len(texts) {
				temp_str += parse(texts[tranID:len(texts)])
				extractRegex(temp_str, &transaction)
				transactions = append(transactions, transaction)
			}
			if !TransactionNum.MatchString(sentence.S) {
				continue
			}

			if tranID != 0 {
				temp_str += parse(texts[tranID : i-(len(sentence.S)-1)])
				tranID = i + 1
				extractRegex(temp_str, &transaction)
				transactions = append(transactions, transaction)
				temp_str = ""
				transaction = Transaction{}
			}
			transaction.Id = TransactionNum.FindString(sentence.S)
			tranID = i + 1
			sentence.S = ""
		}
	}
	//csv(transactions)
	return transactions
}

func parse(array []pdf.Text) string {
	var str string
	var sentence pdf.Text
	var temp pdf.Text

	for _, text := range array {
		if text.Y == sentence.Y {
			sentence.S = text.S
		} else if sentence.X == text.X {
			temp = text
			temp.S = " " + temp.S
			sentence = temp
		} else {
			sentence = text
			sentence.S = " " + sentence.S
		}
		str += sentence.S
	}
	return str
}

func extractRegex(str string, transaction *Transaction) {
	var priceIndex, dateIndex, timeIndex, typeIndex []int

	if TransactionType.MatchString(str) {
		typeIndex = TransactionType.FindStringIndex(str)
		find := str[typeIndex[0]:(typeIndex[1] - 2)]
		find = strings.Trim(find, " ")
		transaction.Typeof = find
	}
	if PriceWithAcronym.MatchString(str) {
		priceIndex = PriceWithAcronym.FindStringIndex(str)
		price_acronym := str[priceIndex[0]:priceIndex[1]]

		transaction.Price = Price.FindString(price_acronym)
		transaction.Acronym = Acronym.FindString(price_acronym)
		transaction.Description = str[priceIndex[1]:len(str)]
	}
	if len(typeIndex) != 0 && len(priceIndex) != 0 {
		find := str[typeIndex[1]-2 : priceIndex[0]]
		find = strings.Trim(find, " ")
		transaction.Status = find
	}
	str = strings.ReplaceAll(str, " ", "")

	if Date.MatchString(str) {
		dateIndex = Date.FindStringIndex(str)
		transaction.Date = str[dateIndex[0]:dateIndex[1]]
	}
	if Time.MatchString(str) {
		timeIndex = Time.FindStringIndex(str)
		transaction.Time = str[timeIndex[0]:timeIndex[1]]
	}

}

func csv(t []Transaction) {
	path := "./alfa.csv"
	file, err := os.Create(path)
	check(err)
	str := "ID;DATE;TIME;TYPE_OF_TRANSACTION;STATUS;PRICE;ACRONYM;DESCRIPTION\n"
	for _, transaction := range t {
		temp := []string{transaction.Id, transaction.Date, transaction.Time,
			transaction.Typeof, transaction.Status, transaction.Price, transaction.Acronym, transaction.Description}
		str += strings.Join(temp, ";")
		str += "\n"
	}
	file.WriteString(str)
	file.Close()
}
