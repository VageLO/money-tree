package action

import (
	"database/sql"
	"errors"
	s "main/structs"
	m "main/modal"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RenameNode(text string, n *tview.TreeNode) error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	nodeReference := n.GetReference().(*s.TreeNode)
	id := nodeReference.Reference.Id

	query := `UPDATE Categories SET title = ? WHERE id = ?`

	if _, err = db.Exec(query, text, id); err != nil {
		return err
	}

	n.SetText(text)
	defer db.Close()

	return nil
}

func RemoveCategory(tree *tview.TreeView) error {
	selectedNode := tree.GetCurrentNode()
	if selectedNode == nil {
		return nil
	}

	node := selectedNode.GetReference().(*s.TreeNode)
	id := node.Reference.Id

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	query := `DELETE FROM Categories WHERE id = ? OR parent_id = ?`

	if _, err = db.Exec(query, id, id); err != nil {
		return err
	}

	selectedNode.ClearChildren()
	node.Parent.RemoveChild(selectedNode)
	defer db.Close()

	return nil
}

func AddCategory(newNode *tview.TreeNode, parentNode *tview.TreeNode) error {
	if err := isEmpty(newNode); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

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

	if err != nil {
		return err
	}

	createdId, _ := result.LastInsertId()

	nodeReference.Reference.Id = createdId
	nodeReference.Reference.ParentId.Scan(parentId)
	nodeReference.Reference.Title = title

	newNode.SetReference(nodeReference)

	parentNode.AddChild(newNode)
	defer db.Close()

	return nil
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

func SelecteByCategoryId(id int64, source *s.Source) error {
	query, err := os.ReadFile("./sql/Select_On_Transactions_Where_CategoryID.sql")
	if err != nil {
		return err
	}

	strId := strconv.FormatInt(id, 10)

	request := string(query)
	request = strings.ReplaceAll(request, "?", strId)
	LoadTransactions(request, source)

	return nil
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

