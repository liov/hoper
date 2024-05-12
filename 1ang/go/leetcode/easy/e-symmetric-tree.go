package easy

/*


101. 对称二叉树
简单

给你一个二叉树的根节点 root ， 检查它是否轴对称。
*/

func isSymmetric(root *TreeNode) bool {
	return check(root, root)
}

func check(p, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if q == nil || p == nil {
		return false
	}
	return q.Val == p.Val && check(p.Left, q.Right) && check(p.Right, q.Left)
}
