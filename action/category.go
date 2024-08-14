package action

import (
	"database/sql"
	"errors"
	m "main/modal"
	s "main/structs"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

func RenameNode(text string, node *tview.TreeNode, source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	nodeReference := node.GetReference().(*s.TreeNode)
	id := nodeReference.Reference.Id

	query := `UPDATE Categories SET title = ? WHERE id = ?`

	_, err = db.Exec(query, text, id)
	check(err)

	node.SetText(text)
	defer db.Close()

	source.Pages.RemovePage("Form")
}

func RemoveCategory(source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	tree := source.CategoryTree

	selectedNode := tree.GetCurrentNode()
	if selectedNode == nil {
		return
	}

	node := selectedNode.GetReference().(*s.TreeNode)
	id := node.Reference.Id

	db, err := sql.Open("sqlite3", "./database.db?_foreign_keys=on")
	check(err)

	query := `DELETE FROM Categories WHERE id = ? OR parent_id = ?`

	_, err = db.Exec(query, id, id)
	check(err)

	selectedNode.ClearChildren()
	node.Parent.RemoveChild(selectedNode)

    // Reload transactions
	LoadTransactions(s.Transactions, source)
    
	defer db.Close()
}

func AddCategory(newNode *tview.TreeNode, parentNode *tview.TreeNode, source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)

	err := isEmpty(newNode)
	check(err)

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	title := newNode.GetText()

	parentReference := parentNode.GetReference().(*s.TreeNode)
	nodeReference := newNode.GetReference().(*s.TreeNode)

	var parentId int64
	var result sql.Result

	if parentReference.Reference != nil {
		parentId = parentReference.Reference.Id

		query := `INSERT INTO Categories (title, parent_id) VALUES (?, ?)`
		result, err = db.Exec(query, title, parentId)
	} else {
		query := `INSERT INTO Categories (title) VALUES (?)`
		result, err = db.Exec(query, title)
	}

	check(err)

	createdId, _ := result.LastInsertId()

	nodeReference.Reference.Id = createdId
	nodeReference.Reference.ParentId.Scan(parentId)
	nodeReference.Reference.Title = title
    
    nodeReference.Selected = func() {
        SelecteByCategoryId(createdId, source)
    }

	newNode.SetReference(nodeReference)
	parentNode.AddChild(newNode)
    
    source.App.SetFocus(source.CategoryTree)
	defer db.Close()
}

func SelectCategories(request string, source *s.Source) ([]string, []s.Category, []*s.TreeNode) {
	defer m.ErrorModal(source.Pages, source.Modal)

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	rootCategories, err := db.Query(request)
	check(err)

	var categoryTitles []string
	var categoryTypes []s.Category
	var categoryNodes []*s.TreeNode

	for rootCategories.Next() {
		var c s.Category
		err := rootCategories.Scan(&c.Id, &c.ParentId, &c.Title)
		check(err)

		categoryTitles = append(categoryTitles, c.Title)
		categoryTypes = append(categoryTypes, c)
		categoryNodes = append(categoryNodes, &s.TreeNode{
			Text:      c.Title,
			Expand:    true,
			Reference: &c,
			Children:  []*s.TreeNode{},
			Selected: func() {
				id := c.Id
				SelecteByCategoryId(id, source)
			},
		})
	}

	defer rootCategories.Close()
	defer db.Close()
	return categoryTitles, categoryTypes, categoryNodes
}

func SelecteByCategoryId(id int64, source *s.Source) {
	strId := strconv.FormatInt(id, 10)

	request := strings.ReplaceAll(s.TransactionsWhereCategoryId, "?", strId)
	LoadTransactions(request, source)
}

func isEmpty(n *tview.TreeNode) error {
	if n.GetText() == "" {
		return errors.New("Empty field")
	}
	return nil
}

func AddNode(target *s.TreeNode, parent *tview.TreeNode) *tview.TreeNode {
	node := tview.NewTreeNode(target.Text).
		SetSelectable(target.Expand || target.Selected != nil).
		SetExpanded(target.Expand).
		SetReference(target)
	if target.Expand {
		node.SetColor(tcell.ColorPurple)
	} else if target.Selected != nil {
		node.SetColor(tcell.ColorGreen)
	}
	if parent != nil {
		target.Parent = parent
	}
	for _, child := range target.Children {
		node.AddChild(AddNode(child, node))
	}
	return node
}
