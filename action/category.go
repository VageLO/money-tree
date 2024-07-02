package action

import (
	"database/sql"
	"log"
	"fmt"
	"errors"
	"os"
	"strconv"
	"strings"
	s "main/structs"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)



func RenameNode(text string, n *tview.TreeNode) error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	
	node_reference := n.GetReference().(*s.TreeNode)
	id := node_reference.Reference.Id
	
	query := `UPDATE Categories SET title = ? WHERE id = ?`

	if _, err = db.Exec(query, text, id); err != nil {
		return err
	}
	
	n.SetText(text)
	defer db.Close()
	
	return nil
}

func RemoveCategory(tree *tview.TreeView) error {
	selected_node := tree.GetCurrentNode()
	if selected_node == nil {
		return nil
	}
	
	node := selected_node.GetReference().(*s.TreeNode)
	id := node.Reference.Id
	
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	
	query := `DELETE FROM Categories WHERE id = ? OR parent_id = ?`

	if _, err = db.Exec(query, id, id); err != nil {
		return err
	}
	
	selected_node.ClearChildren()
	node.Parent.RemoveChild(selected_node)
	defer db.Close()
	
	return nil
}

func AddCategory(new_node *tview.TreeNode, parent_node *tview.TreeNode) error {
	if err := isEmpty(new_node); err != nil {
		return err
	}
	
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	
	title := new_node.GetText()
	
	parent_reference := parent_node.GetReference().(*s.TreeNode)
	node_reference := new_node.GetReference().(*s.TreeNode)
	
	var parent_id int64
	var result sql.Result
	
	if parent_reference.Reference != nil {
		parent_id = parent_reference.Reference.Id
		
		query := `INSERT INTO Categories (title, parent_id) VALUES (?, ?)`
		result, err = db.Exec(query, title, parent_id)
	} else {
		query := `INSERT INTO Categories (title) VALUES (?)`
		result, err = db.Exec(query, title)
	}
	
	if err != nil {
		return err
	}
	
	created_id, _ := result.LastInsertId()
	
	node_reference.Reference.Id = created_id
	node_reference.Reference.Parent_id.Scan(parent_id)
	node_reference.Reference.Title = title
	
	new_node.SetReference(node_reference)
	
	parent_node.AddChild(new_node)
	defer db.Close()
	
	return nil
}

func SelectCategories(request string) ([]string, []s.Category, []s.TreeNode, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, nil, nil, err
	}
	
	root_categories, err := db.Query(request)
	if err != nil {
		return nil, nil, nil, err
	}

	var category_titles []string
	var category_types []s.Category
	var category_nodes []s.TreeNode
	
	for root_categories.Next() {
		var c s.Category
		if err := root_categories.Scan(&c.Id, &c.Parent_id, &c.Title); err != nil {
			return nil, nil, nil, err
		}
		category_titles = append(category_titles, c.Title)
		category_types = append(category_types, c)
		category_nodes = append(category_nodes, s.TreeNode {
			Text: c.Title,
			Expand: true,
			Reference: &c,
			Children: []*s.TreeNode{},
			Selected: func() {
				id := c.Id
				SelectedCategory(id)
			},
		})
	}

	defer root_categories.Close()
	defer db.Close()
	return category_titles, category_types, category_nodes, nil
}

func SelectedCategory(id int64) error {
	query, err := os.ReadFile("./sql/Select_On_Transactions_Where_CategoryID.sql")
    if err != nil {
		return err
	}
	
	str_id := strconv.FormatInt(id, 10)

	request := string(query)
	request = strings.ReplaceAll(request, "?", str_id)
	FillTable(request)
	
	return nil
}

func isEmpty(n *tview.TreeNode) error {
	if n.GetText() == "" {
		return errors.New("Empty field")
	}
	return nil
}