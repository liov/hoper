package hard

import (
	"math"
	"math/rand"
	"time"
)

/*
1044. 最长重复子串
给你一个字符串 s ，考虑其所有 重复子串 ：即，s 的连续子串，在 s 中出现 2 次或更多次。这些出现之间可能存在重叠。

返回 任意一个 可能具有最长长度的重复子串。如果 s 不含重复子串，那么答案为 "" 。



示例 1：

输入：s = "banana"
输出："ana"
示例 2：

输入：s = "abcd"
输出：""


提示：

2 <= s.length <= 3 * 104
s 由小写英文字母组成

https://leetcode-cn.com/problems/longest-duplicate-substring/
*/

/*
#### 方法一：二分查找 + Rabin-Karp 字符串编码

**思路及解法**

记 $s$ 的长度为 $n$。这个问题可以分为两步：从 $n - 1$ 到 $1$ 由大至小遍历选取长度 $L$，判断 $s$ 中是否有长度为 $L$ 的重复子串。从大至小遍历的时候，第一次遇到的符合条件的 $L$，即为最大的符合条件的 $L$，记为 $L_1$，重复的子串为 $s_1$。并且对于任意满足 $L_0 \leq L_1$ 的 $L_0$ 也均符合条件，因为 $s_1$ 的所有子串也是 $s$ 的重复子串。而对于任意满足 $L_2 \gt L_1$ 的 $L_2$，则均不符合条件。因此，我们可以用二分查找的方法，来查找 $L_1$。

剩下的任务便是如何高效判断 $s$ 中是否有长度为 $L$ 的重复子串。我们可以使用 Rabin-Karp 算法对固定长度的字符串进行编码。当两个字符串的编码相同时，则这两个字符串也相同。在 $s$ 中 ${n-L+1}$ 个长度为 $L$ 的子串中，有两个子串的编码相同时，则说明存在长度为 $L$ 的重复子串。具体步骤如下：
1. 首先，我们需要对 $s$ 的每个字符进行编码，得到一个数组 $arr$。因为本题中 $s$ 仅包含小写字母，我们可按照 $\texttt{arr[i] = (int)s.charAt(i) - (int)`a'}$，将所有字母编码为 $0-25$ 之间的数字。比如字符串 $\text{``abcde"}$ 可以编码为数组 $[0, 1, 2, 3, 4]$。
2. 我们将子串看成一个 $26$ 进制的数，它对应的 $10$ 进制数就是它的编码。假设此时我们需要求长度为 $3$ 的子串的编码。那么第一个子串 $\text{``abc''}$ 的编码就是 $h_0 = 0 \times 26^2 + 1 \times 26^1 + 2 \times 26^0 = 28$。更一般地，设 $c_i$ 为 $s$ 的第 $i$ 个字符编码后的数字，$a$ $(a\geq26)$ 为编码的进制，那么有 $h_0 = c_0a^{L-1} + c_1a^{L-2} + ... +c_{L-1}a^1 = \sum_{i=0}^{L-1} c_ia^{L-1-i}$。
3. 上一步我们只求了第一个子串 $\text{``abc''}$ 的编码。当我们要求第二个子串 $\text{``bcd''}$ 的编码时，也可以按照上一步的方法求：$h_1 = 1 \times 26^2 + 2 \times 26^1 + 3 \times 26^0 = 731$，但是这样时间复杂度是 $O(L)$。我们可以在 $h_0$ 的基础上，更快地求出它的编码：$h_1 = (h_0 - 0 \times 26^2) \times 26 + 3 \times 26^0 = 731$。更一般的表达式是：$h_1 = (h_0 \times a - c_0 \times a^L) + c_{L+1}$。这样，我们只需要在常数时间内就可以根据上一个子串的编码求出下一个子串的编码。我们用一个哈希表 $\textit{seen}$ 来存储子串的编码。在求子串的编码时，如果某个子串的编码出现过，则表示存在长度为 $L$ 的重复子串，否则，我们将当前的编码放入 $\textit{seen}$ 中。如果所有编码都不重复，则说明不存在长度为 $L$ 的重复子串。
4. 还有一点需要考虑的是，本题中 $a^L$ 会非常大。一般的做法是需要对编码进行取模来防止溢出，模一般选取编码的信息量的平方的数量级。而取模则会带来哈希碰撞。本题中为了避免碰撞，我们使用双哈希，即用两套进制和模的组合，来对字符串进行编码。只有两种编码都相同时，我们才认为字符串相同。
5. 本题要求返回最长重复子串而不是最长重复子串长度。因此，当存在长度为 $L$ 的子串时，我们的判断函数可以返回重复子串的起点。而当不存在时，可以返回 $-1$ 用做区分。

**复杂度分析**

- 时间复杂度：$O(n \log n)$，其中 $n$ 是字符串 $s$ 的长度。二分查找的时间复杂度为 $O(\log n)$，Rabin-Karp 字符串编码的时间复杂度为 $O(n)$。

- 空间复杂度：$O(n)$，其中 $n$ 是字符串 $s$ 的长。$\textit{arr}$ 和 $\textit{seen}$ 各消耗 $O(n)$ 的空间。
*/
/*
那如何进行hash呢？ 我们可以用一个质数 p ，比如 31 当作底数； 将字符串转化为 sub[0]*p^{m-1}+sub[1]*p^{m-2}...+sub[m-1]。 这其实基本上就是 JDK 中对 string 哈希的默认做法。
而这个哈希计算在滑动过程中，我们也不需要每次都重新计算一遍，可以用上一位置的状态转移过来，如果将hash值看成31进制数的话就是所有位数都左移一位，再去头加尾即可，相信很好理解：
hash = hash * prime - prime^m * (s[i-len] - 'a') + (s[i] - 'a')
这就是滑动窗口的威力。

正常来说这个值会很大可能会导致溢出，所以RK算法应该还要对这个数取模，这样会导致hash冲突；不过用 unsigned long long 存储相当于自动取模了

作者：wfnuser
链接：https://leetcode-cn.com/problems/longest-duplicate-substring/solution/wei-rao-li-lun-rabin-karp-er-fen-sou-suo-3c22/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
*/
func randInt(a, b int) int {
	return a + rand.Intn(b-a)
}

func pow(x, n, mod int) int {
	res := 1
	for ; n > 0; n >>= 1 {
		if n&1 > 0 {
			res = res * x % mod
		}
		x = x * x % mod
	}
	return res
}

func check(arr []byte, m, a1, a2, mod1, mod2 int) int {
	aL1, aL2 := pow(a1, m, mod1), pow(a2, m, mod2)
	h1, h2 := 0, 0
	for _, c := range arr[:m] {
		h1 = (h1*a1 + int(c)) % mod1
		h2 = (h2*a2 + int(c)) % mod2
	}
	// 存储一个编码组合是否出现过
	seen := map[[2]int]bool{{h1, h2}: true}
	for start := 1; start <= len(arr)-m; start++ {
		h1 = (h1*a1 - int(arr[start-1])*aL1 + int(arr[start+m-1])) % mod1
		if h1 < 0 {
			h1 += mod1
		}
		h2 = (h2*a2 - int(arr[start-1])*aL2 + int(arr[start+m-1])) % mod2
		if h2 < 0 {
			h2 += mod2
		}
		// 如果重复，则返回重复串的起点
		if seen[[2]int{h1, h2}] {
			return start
		}
		seen[[2]int{h1, h2}] = true
	}
	// 没有重复，则返回 -1
	return -1
}

func longestDupSubstring(s string) string {
	rand.Seed(time.Now().UnixNano())
	// 生成两个进制
	a1, a2 := randInt(26, 100), randInt(26, 100)
	// 生成两个模
	mod1, mod2 := randInt(1e9+7, math.MaxInt32), randInt(1e9+7, math.MaxInt32)
	// 先对所有字符进行编码
	arr := []byte(s)
	for i := range arr {
		arr[i] -= 'a'
	}

	l, r := 1, len(s)-1
	length, start := 0, -1
	for l <= r {
		m := l + (r-l+1)/2
		idx := check(arr, m, a1, a2, mod1, mod2)
		if idx != -1 { // 有重复子串，移动左边界
			l = m + 1
			length = m
			start = idx
		} else { // 无重复子串，移动右边界
			r = m - 1
		}
	}
	if start == -1 {
		return ""
	}
	return s[start : start+length]
}
