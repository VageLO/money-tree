package action

import (
	m "main/modal"
	s "main/structs"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type nodeReference struct {
	path   string
	isDir  bool
	parent *tview.TreeNode
}

type tree struct {
	*tview.TreeView
	rootNode *tview.TreeNode
}

var rootDir, _ = os.UserHomeDir()

func newTree(source *s.Source, pattern, pageName string) *tree {
	defer m.ErrorModal(source.Pages, source.Modal)

	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed).
		SetReference(newNodeReference(rootDir, true, nil))

	tree := &tree{
		TreeView: tview.NewTreeView().
			SetRoot(root).
			SetCurrentNode(root),
		rootNode: root,
	}

	tree.addNode(root, rootDir, pattern)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		tree.expandOrAddNode(node, source, pattern, pageName)
	})

	return tree
}

func (tree *tree) addNode(directoryNode *tview.TreeNode, path, pattern string) {
	files, err := os.ReadDir(path)
	check(err)

	for _, file := range files {
		if pattern != "" && !file.IsDir() && filepath.Ext(file.Name()) != pattern {
			continue
		}
		node := createTreeNode(file.Name(), file.IsDir(), directoryNode)
		directoryNode.AddChild(node)
	}
}

func (tree tree) expandOrAddNode(node *tview.TreeNode, source *s.Source, pattern, pageName string) {
	defer m.ErrorModal(source.Pages, source.Modal)

	reference := node.GetReference()
	if reference == nil {
		return
	}

	nodeReference := reference.(*nodeReference)
	if !nodeReference.isDir && pattern == "" {
		source.Attachments = append(source.Attachments, nodeReference.path)
		source.Pages.RemovePage(pageName)
		return

	} else if !nodeReference.isDir && pattern != "" {
		source.Pages.RemovePage(pageName)
		ImportForm(source, nodeReference.path)
		return
	}

	children := node.GetChildren()
	if len(children) == 0 {
		path := nodeReference.path
		tree.addNode(node, path, pattern)
	} else {
		node.SetExpanded(!node.IsExpanded())
	}
}

func createTreeNode(fileName string, isDir bool, parent *tview.TreeNode) *tview.TreeNode {
	var parentPath string

	if parent == nil {
		parentPath = rootDir
	} else {
		reference, ok := parent.GetReference().(*nodeReference)
		if !ok {
			parentPath = rootDir
		} else {
			parentPath = reference.path
		}
	}

	var color tcell.Color
	if isDir {
		color = tcell.ColorGreen
	} else {
		color = tview.Styles.PrimaryTextColor
	}

	return tview.NewTreeNode(fileName).
		SetReference(
			newNodeReference(
				filepath.Join(parentPath, fileName),
				isDir,
				parent,
			),
		).
		SetSelectable(true).
		SetColor(color)
}

func newNodeReference(path string, isDir bool, parent *tview.TreeNode) *nodeReference {
	return &nodeReference{
		path:   path,
		isDir:  isDir,
		parent: parent,
	}
}
