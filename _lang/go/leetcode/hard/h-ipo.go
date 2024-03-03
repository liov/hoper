package hard

import (
	"sort"
	"test/leetcode"
)

/*
 IPO
假设 力扣（LeetCode）即将开始 IPO 。为了以更高的价格将股票卖给风险投资公司，力扣 希望在 IPO 之前开展一些项目以增加其资本。 由于资源有限，它只能在 IPO 之前完成最多 k 个不同的项目。帮助 力扣 设计完成最多 k 个不同项目后得到最大总资本的方式。

给你 n 个项目。对于每个项目 i ，它都有一个纯利润 profits[i] ，和启动该项目需要的最小资本 capital[i] 。

最初，你的资本为 w 。当你完成一个项目时，你将获得纯利润，且利润将被添加到你的总资本中。

总而言之，从给定项目中选择 最多 k 个不同项目的列表，以 最大化最终资本 ，并输出最终可获得的最多资本。

答案保证在 32 位有符号整数范围内。



示例 1：

输入：k = 2, w = 0, profits = [1,2,3], capital = [0,1,1]
输出：4
解释：
由于你的初始资本为 0，你仅可以从 0 号项目开始。
在完成后，你将获得 1 的利润，你的总资本将变为 1。
此时你可以选择开始 1 号或 2 号项目。
由于你最多可以选择两个项目，所以你需要完成 2 号项目以获得最大的资本。
因此，输出最后最大化的资本，为 0 + 1 + 3 = 4。
示例 2：

输入：k = 3, w = 0, profits = [1,2,3], capital = [0,1,2]
输出：6


提示：

1 <= k <= 105
0 <= w <= 109
n == profits.length
n == capital.length
1 <= n <= 105
0 <= profits[i] <= 104
0 <= capital[i] <= 109

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/ipo
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/

type ProfitNodeIPO struct {
	Profit   int
	Capitals *CapitalListNodeIPO
	Next     *ProfitNodeIPO
}

func (t *ProfitNodeIPO) Put(profit, capital int) *ProfitNodeIPO {
	next := t.Next
	if profit == t.Profit {
		t.Capitals = t.Capitals.Put(capital)
		return t
	}
	profitNodeIPO := &ProfitNodeIPO{Profit: profit, Capitals: &CapitalListNodeIPO{Capital: capital}}
	if profit > t.Profit {
		profitNodeIPO.Next = t
		return profitNodeIPO
	}
	if next == nil || next.Profit < profit {
		t.Next = profitNodeIPO
		profitNodeIPO.Next = next
		return t
	}
	t.Next = t.Next.Put(profit, capital)
	return t
}

type CapitalListNodeIPO struct {
	Capital int
	Next    *CapitalListNodeIPO
}

func (l *CapitalListNodeIPO) Put(capital int) *CapitalListNodeIPO {
	capitalNode := &CapitalListNodeIPO{Capital: capital}
	next := l.Next
	if capital > l.Capital {
		if next == nil || next.Capital > capital {
			l.Next = capitalNode
			capitalNode.Next = next
			return l
		}
		l.Next = l.Next.Put(capital)
		return l
	}

	capitalNode.Next = l
	return capitalNode

}

type ProfitNodeIPO1 struct {
	Profit int
	Next   *ProfitNodeIPO1
}

func (t *ProfitNodeIPO1) Put(profit int) *ProfitNodeIPO1 {
	profitNode := &ProfitNodeIPO1{Profit: profit}
	next := t.Next
	if profit < t.Profit {
		if next == nil || next.Profit < profit {
			t.Next = profitNode
			profitNode.Next = next
			return t
		}
		t.Next = t.Next.Put(profit)
		return t
	}

	profitNode.Next = t
	return profitNode
}

type CapitalListNodeIPO1 struct {
	Capital int
	Profits *ProfitNodeIPO1
	Next    *CapitalListNodeIPO1
}

func (l *CapitalListNodeIPO1) Put(profit, capital int) *CapitalListNodeIPO1 {
	next := l.Next
	if capital == l.Capital {
		l.Profits = l.Profits.Put(profit)
		return l
	}
	capitalNode := &CapitalListNodeIPO1{Capital: capital, Profits: &ProfitNodeIPO1{Profit: profit}}
	if capital < l.Capital {
		capitalNode.Next = l
		return capitalNode
	}
	if next == nil || next.Capital > capital {
		l.Next = capitalNode
		capitalNode.Next = next
		return l
	}
	l.Next = l.Next.Put(profit, capital)
	return l
}

type MaxHeapIPO []*ProfitNodeIPO1

func NewMaxHeapIPOFromArr(arr []*ProfitNodeIPO1) MaxHeapIPO {
	heap := MaxHeapIPO(arr)
	for i := 1; i < len(arr); i++ {
		heap.adjustUp(i)
	}
	return heap
}

func (heap MaxHeapIPO) Put(v *ProfitNodeIPO1) {
	if v.Profit > heap[0].Profit {
		return
	}

	heap[0] = v
	i := 0
	heap.adjustDown(i)
}

func (heap MaxHeapIPO) adjustUp(i int) {
	p := leetcode.parent(i)
	for p >= 0 && heap[i].Profit > heap[p].Profit {
		heap.swap(i, p)
		i = p
		p = leetcode.parent(i)
	}

}

func (heap MaxHeapIPO) swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}
func (heap MaxHeapIPO) adjustDown(i int) {
	child := leetcode.leftChild(i)
	for child < len(heap) {
		if child+1 < len(heap) && heap[child+1].Profit > heap[child].Profit {
			child++
		}
		if heap[i].Profit >= heap[child].Profit {
			break
		}
		heap.swap(i, child)
		i = child
		child = leetcode.leftChild(i)
	}
}

type CapitalListNodeIPO2 struct {
	Capital int
	Profits []int
	Next    *CapitalListNodeIPO2
}

func (node *CapitalListNodeIPO2) Put(capital, profit int) *CapitalListNodeIPO2 {
	next := node.Next
	if capital == node.Capital {
		node.Profits = append(node.Profits, profit)
		return node
	}
	capitalNode := &CapitalListNodeIPO2{Capital: capital, Profits: []int{profit}}
	if capital < node.Capital {
		capitalNode.Next = node
		return capitalNode
	}
	if next == nil || next.Capital > capital {
		node.Next = capitalNode
		capitalNode.Next = next
		return node
	}
	node.Next = node.Next.Put(profit, capital)
	return node
}

type IPO struct {
	profits, capital []int
}

func (ipo *IPO) Len() int {
	return len(ipo.capital)
}

func (ipo *IPO) Less(i, j int) bool {
	return ipo.capital[i] < ipo.capital[j]
}

func (ipo *IPO) Swap(i, j int) {
	ipo.profits[i], ipo.profits[j] = ipo.profits[j], ipo.profits[i]
	ipo.capital[i], ipo.capital[j] = ipo.capital[j], ipo.capital[i]
}

func findMaximizedCapital(k int, w int, profits []int, capital []int) int {

	sort.Sort(&IPO{profits, capital})

	var profitsHeap leetcode.MaxHeap

	var i, j int
	for {
		for i < len(capital) && capital[i] <= w {
			profitsHeap = append(profitsHeap, profits[i])
			profitsHeap.adjustUp(len(profitsHeap) - 1)
			i++
		}
		if len(profitsHeap) == 0 || profitsHeap[0] == 0 {
			break
		}
		w += profitsHeap[0]
		j++
		if j == k {
			return w
		}
		profitsHeap[0] = 0
		profitsHeap.adjustDown(0)
	}
	return w
}
