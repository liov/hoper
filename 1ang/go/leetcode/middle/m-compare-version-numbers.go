package middle

/*
165. 比较版本号

给你两个版本号 version1 和 version2 ，请你比较它们。

版本号由一个或多个修订号组成，各修订号由一个 '.' 连接。每个修订号由 多位数字 组成，可能包含 前导零 。每个版本号至少包含一个字符。修订号从左到右编号，下标从 0 开始，最左边的修订号下标为 0 ，下一个修订号下标为 1 ，以此类推。例如，2.5.33 和 0.1 都是有效的版本号。

比较版本号时，请按从左到右的顺序依次比较它们的修订号。比较修订号时，只需比较 忽略任何前导零后的整数值 。也就是说，修订号 1 和修订号 001 相等 。如果版本号没有指定某个下标处的修订号，则该修订号视为 0 。例如，版本 1.0 小于版本 1.1 ，因为它们下标为 0 的修订号相同，而下标为 1 的修订号分别为 0 和 1 ，0 < 1 。

返回规则如下：

如果 version1 > version2 返回 1，
如果 version1 < version2 返回 -1，
除此之外返回 0。

示例 1：

输入：version1 = "1.01", version2 = "1.001"
输出：0
解释：忽略前导零，"01" 和 "001" 都表示相同的整数 "1"
示例 2：

输入：version1 = "1.0", version2 = "1.0.0"
输出：0
解释：version1 没有指定下标为 2 的修订号，即视为 "0"
示例 3：

输入：version1 = "0.1", version2 = "1.1"
输出：-1
解释：version1 中下标为 0 的修订号是 "0"，version2 中下标为 0 的修订号是 "1" 。0 < 1，所以 version1 < version2
示例 4：

输入：version1 = "1.0.1", version2 = "1"
输出：1
示例 5：

输入：version1 = "7.5.2.4", version2 = "7.5.3"
输出：-1

提示：

1 <= version1.length, version2.length <= 500
version1 和 version2 仅包含数字和 '.'
version1 和 version2 都是 有效版本号
version1 和 version2 的所有修订号都可以存储在 32 位整数 中

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/compare-version-numbers
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
func compareVersion(version1 string, version2 string) int {
	for {
		part1, idx1, part2, idx2 := getEffectPart(version1, version2)
		if len(part1) > len(part2) {
			return 1
		}
		if len(part1) < len(part2) {
			return -1
		}
		for i := range part1 {
			if part1[i] > part2[i] {
				return 1
			}
			if part1[i] < part2[i] {
				return -1
			}
		}
		if len(version1) == idx1 && len(version2) == idx2 {
			return 0
		}
		if idx1 == len(version1) {
			version1 = ""
		} else {
			version1 = version1[idx1+1:]
		}
		if idx2 == len(version2) {
			version2 = ""
		} else {
			version2 = version2[idx2+1:]
		}
	}
}

func getEffectPart(version1 string, version2 string) (version1Part string, effectEnd1 int, version2Part string, effectEnd2 int) {
	effectStart1 := -1
	effectStart2 := -1

	if version1 != "" {
		for i := range version1 {
			if effectStart1 == -1 && version1[i] != '0' && version1[i] != '.' {
				effectStart1 = i
			}
			if effectStart2 == -1 && effectEnd2 == 0 && i < len(version2) && version2[i] != '0' && version2[i] != '.' {
				effectStart2 = i
			}
			if effectEnd1 == 0 && version1[i] == '.' {
				effectEnd1 = i
			}
			if effectEnd2 == 0 && i < len(version2) && version2[i] == '.' {
				effectEnd2 = i
			}
			if effectEnd1 != 0 {
				break
			}
		}
		if effectEnd1 == 0 {
			effectEnd1 = len(version1)
		}
		if effectStart1 == -1 {
			version1Part = "0"
		} else {
			version1Part = version1[effectStart1:effectEnd1]
		}
	} else {
		version1Part = "0"
	}

	if version2 != "" {
		if effectEnd1 >= len(version2) {
			version2Part = version2
			effectEnd2 = len(version2)
		}
		if effectEnd2 == 0 {
			for i := effectEnd1; i < len(version2); i++ {
				if effectEnd2 == 0 && version2[i] == '.' {
					effectEnd2 = i
				}
				if effectStart2 == -1 && effectEnd2 == 0 && version2[i] != '0' && version2[i] != '.' {
					effectStart2 = i
				}
			}
		}
		if effectEnd2 == 0 {
			effectEnd2 = len(version2)
		}
		if effectStart2 == -1 {
			version2Part = "0"
		} else {
			version2Part = version2[effectStart2:effectEnd2]
		}
	} else {
		version2Part = "0"
	}
	return
}

func compareVersion2(version1 string, version2 string) int {
	for {
		part1, idx1 := getEffectPart2(version1)
		part2, idx2 := getEffectPart2(version2)
		if len(part1) > len(part2) {
			return 1
		}
		if len(part1) < len(part2) {
			return -1
		}
		for i := range part1 {
			if part1[i] > part2[i] {
				return 1
			}
			if part1[i] < part2[i] {
				return -1
			}
		}
		if len(version1) == idx1 && len(version2) == idx2 {
			return 0
		}
		if idx1 == len(version1) {
			version1 = ""
		} else {
			version1 = version1[idx1+1:]
		}
		if idx2 == len(version2) {
			version2 = ""
		} else {
			version2 = version2[idx2+1:]
		}
	}
}

func getEffectPart2(version string) (versionPart string, effectEnd int) {
	if version == "" {
		return "0", 0
	}
	effectStart := -1
	for i, c := range version {
		if c == '.' {
			effectEnd = i
			break
		}
		if effectStart == -1 && c != '0' {
			effectStart = i
		}
	}
	if effectEnd == 0 {
		effectEnd = len(version)
	}
	if effectStart == -1 {
		return "0", effectEnd
	}
	return version[effectStart:effectEnd], effectEnd
}
