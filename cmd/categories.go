package cmd

import (
	"fmt"
	"main/action"
	m "main/modal"
	s "main/structs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	add func(target *s.TreeNode, parent *tview.TreeNode) *tview.TreeNode
)

func MakeTree() *s.TreeNode {
	defer m.ErrorModal(source.Pages, source.Modal)

	_, _, category_nodes, err := action.SelectCategories(`SELECT * FROM Categories WHERE parent_id IS NULL`)
	check(err)

	for i, node := range category_nodes {
		query := fmt.Sprintf(`SELECT * FROM Categories WHERE parent_id = %v`, node.Reference.Id)
		_, _, children_nodes, err := action.SelectCategories(query)
		check(err)
		category_nodes[i].Children = children_nodes
	}

	var rootNode = &s.TreeNode{
		Text:     ".",
		Children: category_nodes,
	}
	return rootNode
}

func CategoryTree() {
	rootNode := MakeTree()
	source.CategoryTree.SetBorder(true).
		SetTitle("Category Tree")

	// Add nodes
	add = func(target *s.TreeNode, parent *tview.TreeNode) *tview.TreeNode {
		node := tview.NewTreeNode(target.Text).
			SetSelectable(target.Expand || target.Selected != nil).
			SetExpanded(target == rootNode).
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
			node.AddChild(add(child, node))
		}
		return node
	}
	root := add(rootNode, nil)
	source.CategoryTree.SetRoot(root).
		SetCurrentNode(root).
		SetSelectedFunc(func(n *tview.TreeNode) {
			original := n.GetReference().(*s.TreeNode)
			if original.Expand {
				n.SetExpanded(!n.IsExpanded())
			}
			if original.Selected != nil {
				original.Selected()
			}
		})

	source.CategoryTree.GetRoot().ExpandAll()
}

