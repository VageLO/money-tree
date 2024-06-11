package cmd

import (
	"database/sql"
	"log"

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
}

type category_type struct {
	id        int
	title     string
}

var (
	tree = tview.NewTreeView().SetAlign(false).SetTopLevel(1).SetGraphics(true).SetPrefixes(nil)
	add  func(target *node, parent *tview.TreeNode) *tview.TreeNode
)

func MakeTree() *node {

	var rootNode = &node{
		text: ".",
		children: []*node{
			{text: "Category 1", expand: true},
			{text: "Category 2", expand: true},
			{text: "Category 3", expand: true, children: []*node{
				{text: "Child node", expand: true},
				{text: "Child node", expand: true},
				{text: "Selected child node", selected: func() {
					// Updating table on selected node
					//FillTable()
					table.SetBorder(true).SetTitle("Categories")
				}, expand: true, children: []*node{{text: "test", expand: true}}},
			}},
		}}
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

func RenameNode() {
	node := tree.GetCurrentNode()
	FillTreeAndListForm(node, nil)
	pages.AddPage("Dialog", Dialog(form), true, true)
}

func RemoveNode() {
	selected_node := tree.GetCurrentNode()
	if selected_node == nil {
		return
	}
	selected_node.ClearChildren()
	reference := selected_node.GetReference().(*node)
	reference.parent.RemoveChild(selected_node)
}

func AddNode() {

	root := tree.GetRoot()

	n := &node{
		text:   "",
		expand: true,
	}
	new_node := add(n, root)

	FillTreeAndListForm(new_node, nil)

	var selected_dropdown *tview.TreeNode
	var options []string
	options = append(options, root.GetText())

	for _, children := range root.GetChildren() {
		options = append(options, children.GetText())
	}

	initial := 0

	selected_node := tree.GetCurrentNode()
	if selected_node != nil {
		for idx, title := range options {
			if title == selected_node.GetText() {
				initial = idx
			}
		}
	}

	form.AddDropDown("Categories", options, initial, func(option string, optionIndex int) {
		for _, children := range root.GetChildren() {
			if children.GetText() == option {
				selected_dropdown = children
				reference := new_node.GetReference().(*node)
				reference.parent = children
				new_node.SetReference(reference)
			}
		}
		if root.GetText() == option {
			selected_dropdown = root
		}
	})

	form.AddButton("Add", func() {
		selected_dropdown.AddChild(new_node)
		pages.RemovePage("Dialog")
	})

	pages.AddPage("Dialog", Dialog(form), true, true)
}

func SelectCategories() ([]string, []category_type) {
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	root_categories, err := db.Query(`SELECT id, title FROM Categories`)
	check(err)

	var category_titles []string
	var category_types []category_type
	
	for root_categories.Next() {
		var c category_type
		if err := root_categories.Scan(&c.id, &c.title); err != nil {
			log.Fatal(err)
		}
		category_titles = append(category_titles, c.title)
		category_types = append(category_types, c)
	}

	defer root_categories.Close()
	db.Close()
	return category_titles, category_types
}
