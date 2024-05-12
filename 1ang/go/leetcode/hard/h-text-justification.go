package hard

import (
	"bytes"
)

/*
68. 文本左右对齐

给定一个单词数组和一个长度 maxWidth，重新排版单词，使其成为每行恰好有 maxWidth 个字符，且左右两端对齐的文本。

你应该使用“贪心算法”来放置给定的单词；也就是说，尽可能多地往每行中放置单词。必要时可用空格 ' ' 填充，使得每行恰好有 maxWidth 个字符。

要求尽可能均匀分配单词间的空格数量。如果某一行单词间的空格不能均匀分配，则左侧放置的空格数要多于右侧的空格数。

文本的最后一行应为左对齐，且单词之间不插入额外的空格。

说明:

单词是指由非空格字符组成的字符序列。
每个单词的长度大于 0，小于等于 maxWidth。
输入单词数组 words 至少包含一个单词。
示例:

输入:
words = ["This", "is", "an", "example", "of", "text", "justification."]
maxWidth = 16
输出:
[
   "This    is    an",
   "example  of text",
   "justification.  "
]
示例 2:

输入:
words = ["What","must","be","acknowledgment","shall","be"]
maxWidth = 16
输出:
[
  "What   must   be",
  "acknowledgment  ",
  "shall be        "
]
解释: 注意最后一行的格式应为 "shall be    " 而不是 "shall     be",
     因为最后一行应为左对齐，而不是左右两端对齐。
     第二行同样为左对齐，这是因为这行只包含一个单词。
示例 3:

输入:
words = ["Science","is","what","we","understand","well","enough","to","explain",
         "to","a","computer.","Art","is","everything","else","we","do"]
maxWidth = 20
输出:
[
  "Science  is  what we",
  "understand      well",
  "enough to explain to",
  "a  computer.  Art is",
  "everything  else  we",
  "do                  "
]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/text-justification
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/

func fullJustify(words []string, maxWidth int) []string {
	var ret []string
	row := bytes.NewBuffer([]byte{})
	width := 0
	i := 0
	for len(words) > 0 {

		if width+len(words[i])+1 > maxWidth || i == len(words)-1 {
			if i == len(words)-1 && width+len(words[i]) < maxWidth {
				endSpaceNum := maxWidth - width - len(words[i])
				for j := 0; j <= i; j++ {
					row.WriteString(words[j])
					if j < i {
						row.WriteByte(' ')
					}
					if j == i {
						for x := 0; x < endSpaceNum; x++ {
							row.WriteByte(' ')
						}
					}
				}
				ret = append(ret, row.String())
				return ret
			}

			if width+len(words[i]) == maxWidth {
				width += len(words[i]) + 1
				i++
			}

			spaceNum := maxWidth - width + i
			arr := []int{spaceNum}
			if i > 1 {
				spaceNum = (maxWidth-width+1)/(i-1) + 1
				arr = make([]int, i-1)
				for x := 0; x < i-1; x++ {
					arr[x] = spaceNum
				}
				mod := (maxWidth - width + 1) % (i - 1)
				y := 2
				for i-y != 0 {
					add := mod / (i - y)
					mod = mod % (i - y)
					if add != 0 {
						for x := 0; x < i-y; x++ {
							arr[x] += add
						}
					} else {
						break
					}
					y++
				}
				for x := 0; x < mod; x++ {
					arr[x] += 1
				}
			}
			for j := 0; j < i; j++ {
				row.WriteString(words[j])
				if j != 0 && j == i-1 {
					break
				}
				for x := 0; x < arr[j]; x++ {
					row.WriteByte(' ')
				}
			}
			ret = append(ret, row.String())
			row.Reset()
			words = words[i:]
			width = 0
			i = 0

			continue
		}
		width += len(words[i]) + 1
		i++
	}
	return ret
}
