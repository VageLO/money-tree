package cmd

import (
	"log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type node struct {
	text     string
	expand   bool
	selected func()
	children []*node
	parent *tview.TreeNode
}

var (
	tree = tview.NewTreeView().SetAlign(false).SetTopLevel(1).SetGraphics(true).SetPrefixes(nil)
)

func MakeTree() *node{
	tableData := `OrderDate|Region|Rep|Item|Units|UnitCost|Total
1/6/2017|East|Jones|Pencil|95|1.99|189.05
1/23/2017|Central|Kivell|Binder|50|19.99|999.50`
		
	var rootNode = &node{
	text: ".",
	children: []*node{
		{text: "Expand all", selected: func() { tree.GetRoot().ExpandAll() }},
		{text: "Collapse all", selected: func() {
			for _, child := range tree.GetRoot().GetChildren() {
				child.CollapseAll()
			}
		}},
		{text: "Root node", expand: true, children: []*node{
			{text: "Child node", expand: true},
			{text: "Child node", expand: true},
			{text: "Selected child node", selected: func() {
				// Updating table on selected node
				FillTable(tableData)
				table.SetBorder(true).SetTitle("Categories")
			}, expand: true, children: []*node {{text: "test", expand: true}}},
		}},
	}}
	return rootNode
}

func TreeView() *tview.TreeView {
	rootNode := MakeTree()
	tree.SetBorder(true).
		SetTitle("Category Tree")

	// Add nodes
	var add func(target *node, parent *tview.TreeNode) *tview.TreeNode
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
	node := tview.NewTreeNode("")
	FillTreeAndListForm(node, nil)
	
	root := tree.GetRoot()
	var sel *tview.TreeNode
	var options []string
	
	for _, children := range root.GetChildren() {
		options = append(options, children.GetText())
	}
	
	form.AddDropDown("Categories", options, 0, func(option string, optionIndex int) {
		for _, children := range root.GetChildren() {
			log.Println(option)
			if children.GetText() == option {
				sel = children
			}
		}
	})
	form.AddButton("Add", func() {
		sel.AddChild(node)
	})
	pages.AddPage("Dialog", Dialog(form), true, true)
}