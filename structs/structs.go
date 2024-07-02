package structs

import (
	"database/sql"
	"errors"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

type Transaction struct {
	id               int64
	account_id       int64
	to_account_id    sql.NullInt64
	category_id      int64
	transaction_type string
	date             string
	amount           float64
	to_amount        sql.NullFloat64
	account          string
	to_account       sql.NullString
	category         string
	description      string
}

type Account struct {
	Id       int64
	Title    string
	Currency string
	Balance  float64
}

type Category struct {
	Id        int64
	Parent_id sql.NullInt64
	Title     string
}

type Cell struct {
	Row        int
	Column     int
	Text       string
	Selectable bool
	Color      tcell.Color
	Reference  interface{}
}

type TreeNode struct {
	Text      string
	Expand    bool
	Selected  func()
	Children  []*TreeNode
	Parent    *tview.TreeNode
	Reference *Category
}

func (a *Account) isEmpty() error {
	if a.Title == "" || a.Currency == "" {
		return errors.New("Empty field or can't be zero")
	}
	return nil
}
