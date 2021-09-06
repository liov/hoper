package leetcode

/*
设计一个算法，找出数组中最小的k个数。以任意顺序返回这k个数均可。

示例：

输入： arr = [1,3,5,7,2,4,6,8], k = 4
输出： [1,2,3,4]
提示：

0 <= len(arr) <= 100000
0 <= k <= min(100000, len(arr))

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/smallest-k-lcci
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/

type MaxHeap []int

func New(arr []int, k int) MaxHeap {
	for i := 1; i < k; i++ {
		adjustUp(arr, i)
	}
	maxHeap := make(MaxHeap, k)
	copy(maxHeap, arr[:k])
	return maxHeap
}
func (m MaxHeap) Put(v int) {
	if v > m[0] {
		return
	}

}

func Swap(arr []int, i, j int) { arr[i], arr[j] = arr[j], arr[i] }
func parent(pos int) int       { return (pos - 1) / 2 }
func leftChild(pos int) int    { return pos*2 + 1 }

func adjustUp(heap []int, pos int) {

	//比父节点大就上移
	var p = parent(pos)
	for p >= 0 && heap[pos] > heap[p] {
		Swap(heap, p, pos)
		pos = p
		p = parent(pos)
	}
}

func adjustDown(heap []int, pos, k int) {
	if heap[pos] >= heap[0] {
		return
	}
	Swap(heap, pos, 0)

	p := 0
	var child = leftChild(p)
	for child < k {
		if child+1 < k && heap[child+1] > heap[child] {
			child++
		}
		if heap[p] >= heap[child] {
			break
		}
		Swap(heap, p, child)
		p = child
		child = leftChild(p)
	}

}

func smallestK(arr []int, k int) []int {
	for i := 1; i < k; i++ {
		adjustUp(arr, i)
	}
	for i := k; i < len(arr); i++ {
		adjustDown(arr, i, k)
	}
	return arr[0:k]
}
