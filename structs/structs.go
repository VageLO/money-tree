package structs

import (
	"database/sql"
	"errors"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Transaction struct {
	Id               int64
	AccountId        int64
	ToAccountId      sql.NullInt64
	CategoryId       int64
	TransactionType string
	Date             string
	Amount           float64
	ToAmount         sql.NullFloat64
	ToAccount        sql.NullString
	Account          string
	Category         string
	Description      string
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
	Columns []string
	Index int
	Data []string
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
}

func (a *Account) isEmpty() error {
	if a.Title == "" || a.Currency == "" {
		return errors.New("Empty field or can't be zero")
	}
	return nil
}
