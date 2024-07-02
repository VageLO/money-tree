package cmd

import (
	"fmt"
	m "main/modal"
	"main/action"
	
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	tree = tview.NewTreeView().SetAlign(false).SetTopLevel(1).SetGraphics(true).SetPrefixes(nil)
	add  func(target *TreeNode, parent *tview.TreeNode) *tview.TreeNode
)

func MakeTree() *TreeNode {
	defer m.ErrorModal(pages, modal)
	
	_, _, category_nodes, err := action.SelectCategories(`SELECT * FROM Categories WHERE parent_id IS NULL`)
	check(err)
	
	for i, node := range category_nodes {
		query := fmt.Sprintf(`SELECT * FROM Categories WHERE parent_id = %v`, node.reference.id)
		_, _, children_nodes, err := action.SelectCategories(query)
		check(err)
		category_nodes[i].children = children_nodes
	}
	
	var rootNode = &TreeNode{
		text: ".",
		children: category_nodes,
	}
	return rootNode
}

func TreeView() {
	rootNode := MakeTree()
	tree.SetBorder(true).
		SetTitle("Category Tree")

	// Add nodes
	add = func(target *TreeNode, parent *tview.TreeNode) *tview.TreeNode {
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
			original := n.GetReference().(*TreeNode)
			if original.expand {
				n.SetExpanded(!n.IsExpanded())
			}
			if original.selected != nil {
				original.selected()
			}
		})

	tree.GetRoot().ExpandAll()
}