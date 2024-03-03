package middle

/*
437. 路径总和 III
给定一个二叉树的根节点 root ，和一个整数 targetSum ，求该二叉树里节点值之和等于 targetSum 的 路径 的数目。

路径 不需要从根节点开始，也不需要在叶子节点结束，但是路径方向必须是向下的（只能从父节点到子节点）。

示例 1：

输入：root = [10,5,-3,3,2,null,11,3,-2,null,1], targetSum = 8
输出：3
解释：和等于 8 的路径有 3 条，如图所示。
示例 2：

输入：root = [5,4,8,11,null,13,4,7,2,null,null,5,1], targetSum = 22
输出：3

提示:

二叉树的节点个数的范围是 [0,1000]
-10^9 <= Node.val <= 10^9
-1000 <= targetSum <= 1000

https://leetcode-cn.com/problems/path-sum-iii/
*/
func pathSum(root *TreeNode, targetSum int) int {
	var result int
	preorderTraversal(root, nil, targetSum, &result)
	return result
}

// 前序
func preorderTraversal(root *TreeNode, sum []int, targetSum int, result *int) {
	if root == nil {
		return
	}
	sum = append(sum, root.Val)
	if root.Val == targetSum {
		*result++
	}
	cpy := searchBFS(sum, root.Val, targetSum, result)
	preorderTraversal(root.Left, sum, targetSum, result)
	preorderTraversal(root.Right, cpy, targetSum, result)
}

func searchBFS(sum []int, val, targetSum int, result *int) []int {
	for i := 0; i < len(sum)-1; i++ {
		sum[i] += val
		if sum[i] == targetSum {
			*result++
		}
	}
	cpy := make([]int, len(sum))
	copy(cpy, sum)
	return cpy
}
