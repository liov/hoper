package middle

import "fmt"

// val > 0
type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

func NewTree(arr []int) *TreeNode {
	root := &TreeNode{Val: arr[0]}
	newTree(arr, root, 1, 2)
	return root
}

func newTree(arr []int, root *TreeNode, left, right int) {
	if left < len(arr) && arr[left] != 0 {
		root.Left = &TreeNode{Val: arr[left]}
		newTree(arr, root.Left, leftChild(left), rightChild(left))
	}
	if right < len(arr) && arr[right] != 0 {
		root.Right = &TreeNode{Val: arr[right]}
		newTree(arr, root.Right, leftChild(right), rightChild(right))
	}
}

type MaxBinaryTreeNode struct {
	Val         int
	Left, Right *MaxBinaryTreeNode
}

// 中序
func InorderTraversal(root *TreeNode) {
	if root == nil {
		return
	}
	InorderTraversal(root.Left)
	fmt.Println(root)
	InorderTraversal(root.Right)
}

// 前序
func PreorderTraversal(root *TreeNode) {
	if root == nil {
		return
	}
	fmt.Println(root)
	PreorderTraversal(root.Left)
	PreorderTraversal(root.Right)
}

// 后序
func PostorderTraversal(root *TreeNode) {
	if root == nil {
		return
	}
	PostorderTraversal(root.Left)
	PostorderTraversal(root.Right)
	fmt.Println(root)
}
