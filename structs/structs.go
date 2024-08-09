package structs

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Transaction struct {
	Id              int64
	AccountId       int64
	ToAccountId     sql.NullInt64
	CategoryId      int64
	TransactionType string
	Date            string
	Amount          float64
	ToAmount        sql.NullFloat64
	ToAccount       sql.NullString
	Account         string
	Category        string
	Description     string
}

type Account struct {
	Id       int64
	Title    string
	Currency string
	Balance  float64
}

type Category struct {
	Id       int64
	ParentId sql.NullInt64
	Title    string
}

type Cell struct {
	Row        int
	Column     int
	Text       string
	Selectable bool
	Color      tcell.Color
	Reference  interface{}
}

type Row struct {
	Columns     []string
	Index       int
	Data        []string
	Transaction Transaction
}

type TreeNode struct {
	Text      string
	Expand    bool
	Selected  func()
	Children  []*TreeNode
	Parent    *tview.TreeNode
	Reference *Category
}

type Source struct {
	App          *tview.Application
	AccountList  *tview.List
	CategoryTree *tview.TreeView
	Form         *tview.Form
	Table        *tview.Table
	FileTable    *tview.Table
	Modal        *tview.Modal
	Pages        *tview.Pages
	Columns      []string
	Attachments  []string
	Imports      []string
}

type Statistics struct {
	Debit int64
	Credit int64
	Category string
}

func (a *Account) IsEmpty() error {
	if a.Title == "" || a.Currency == "" {
		return errors.New("Fields [Title, Currency] can't be empty")
	}
	a.Title = strings.TrimSpace(a.Title)
	a.Currency = strings.TrimSpace(a.Currency)
	return nil
}

func (t *Transaction) IsEmpty() error {
	if t.Date == "" || t.Amount <= 0 {
		if t.Amount <= 0 {
			return errors.New("Amount can't be negative or zero!")
		}
		//test := fmt.Sprintf("%+v", t)
		return errors.New("Fields [Date, Amount] must be filled")
	} else if t.AccountId == t.ToAccountId.Int64 {
		return errors.New("Select different accounts")
	}

	_, err := time.Parse("2006-01-02", t.Date)
	if err != nil {
		return errors.New("Allowed date format (YYYY-MM-DD)")
	}

	t.Description = strings.TrimSpace(t.Description)
	return nil
}
