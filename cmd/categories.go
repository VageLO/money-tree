package cmd

import (
	"database/sql"
	"log"
	"fmt"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

type node struct {
	text     string
	expand   bool
	selected func()
	children []*node
	parent   *tview.TreeNode
	reference *category_type
}

type category_type struct {
	id        int64
	parent_id sql.NullInt64
	title     string
}

var (
	tree = tview.NewTreeView().SetAlign(false).SetTopLevel(1).SetGraphics(true).SetPrefixes(nil)
	add  func(target *node, parent *tview.TreeNode) *tview.TreeNode
)

func MakeTree() *node {
	
	_, _, category_nodes := SelectCategories(`SELECT * FROM Categories WHERE parent_id IS NULL`)
	
	for i, node := range category_nodes {
		query := fmt.Sprintf(`SELECT * FROM Categories WHERE parent_id = %v`, node.reference.id)
		_, _, children_nodes := SelectCategories(query)
		category_nodes[i].children = children_nodes
	}
	
	var rootNode = &node{
		text: ".",
		children: category_nodes,
	}
	return rootNode
}

func TreeView() *tview.TreeView {
	rootNode := MakeTree()
	tree.SetBorder(true).
		SetTitle("Category Tree")

	// Add nodes
	add = func(target *node, parent *tview.TreeNode) *tview.TreeNode {
		node := tview.NewTreeNode(target.text).
			SetSelectable(target.expand || target.selected != nil).
			SetExpanded(target == rootNode).
			SetReference(target)
		if target.expand {
			node.SetColor(tcell.ColorPurple)
		} else if target.selected != nil {
			node.SetColor(tcell.ColorGreen)
		}
		if parent != nil {
			target.parent = parent
		}
		for _, child := range target.children {
			node.AddChild(add(child, node))
		}
		return node
	}
	root := add(rootNode, nil)
	tree.SetRoot(root).
		SetCurrentNode(root).
		SetSelectedFunc(func(n *tview.TreeNode) {
			original := n.GetReference().(*node)
			if original.expand {
				n.SetExpanded(!n.IsExpanded())
			} else if original.selected != nil {
				original.selected()
			}
		})

	tree.GetRoot().ExpandAll()

	return tree
}

func RenameNode(text string, n *tview.TreeNode) {
	defer ErrorModal()
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	node_reference := n.GetReference().(*node)
	id := node_reference.reference.id
	
	query := `UPDATE Categories SET title = ? WHERE id = ?`

	_, err = db.Exec(query, text, id)
	check(err)
	
	n.SetText(text)
	defer db.Close()
}

func RemoveNode() {
	selected_node := tree.GetCurrentNode()
	if selected_node == nil {
		return
	}
	defer ErrorModal()
	
	node := selected_node.GetReference().(*node)
	id := node.reference.id
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	query := `DELETE FROM Categories WHERE id = ? OR parent_id = ?`

	_, err = db.Exec(query, id, id)
	check(err)
	
	selected_node.ClearChildren()
	node.parent.RemoveChild(selected_node)
	defer db.Close()
}

func AddCategory(new_node *tview.TreeNode, parent_node *tview.TreeNode) {
	defer ErrorModal()
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	title := new_node.GetText()
	
	parent_reference := parent_node.GetReference().(*node)
	node_reference := new_node.GetReference().(*node)
	
	var parent_id int64
	var result sql.Result
	
	if parent_reference.reference != nil {
		parent_id = parent_reference.reference.id
		query := `INSERT INTO Categories (title, parent_id) VALUES (?, ?)`

		result, err = db.Exec(query, title, parent_id)
		check(err)
	} else {
		query := `INSERT INTO Categories (title) VALUES (?)`

		result, err = db.Exec(query, title)
		check(err)
	}
	
	created_id, _ := result.LastInsertId()
	
	node_reference.reference.id = created_id
	// TODO: Make sure the correct data is stored
	node_reference.reference.parent_id.Scan(parent_id)
	node_reference.reference.title = title
	
	new_node.SetReference(node_reference)
	
	parent_node.AddChild(new_node)
	defer db.Close()
}

func SelectCategories(request string) ([]string, []category_type, []*node) {
	defer ErrorModal()
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	root_categories, err := db.Query(request)
	check(err)

	var category_titles []string
	var category_types []category_type
	var category_nodes []*node
	
	for root_categories.Next() {
		var c category_type
		if err := root_categories.Scan(&c.id, &c.parent_id, &c.title); err != nil {
			log.Fatal(err)
		}
		category_titles = append(category_titles, c.title)
		category_types = append(category_types, c)
		category_nodes = append(category_nodes, &node {
			text: c.title,
			expand: true,
			reference: &c,
			children: []*node{},
		})
	}

	defer root_categories.Close()
	defer db.Close()
	return category_titles, category_types, category_nodes
}
