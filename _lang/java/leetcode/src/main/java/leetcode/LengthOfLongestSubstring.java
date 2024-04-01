package leetcode;
import java.lang.Math;

public class LengthOfLongestSubstring {
  public static int lengthOfLongestSubstring(String s) {
    var len = s.length();
    if (len <= 1) {
      return len;
    }
    var maxLength = 0;
    var left = -1;
    var intArray = new int[128];
    var idx = 0;
    for (var i = 0; i < len; i++) {
      idx = s.charAt(i);
      left = Math.max(left, intArray[idx] - 1);
      intArray[idx] = i + 1;
      maxLength = Math.max(maxLength, i - left);
    }
    return maxLength;
  }

  public static int lengthOfLongestSubstringV2(String s) {
    var len = s.length();
    if (len <= 1) {
      return len;
    }
    var maxLength = 0;
    var currentLength = 0;
    var intArray = new int[128];
    var idx = 0;
    for (var i = 0; i < len; i++) {
      idx = s.charAt(i);
      if (intArray[idx] != 0 && intArray[idx] >= i - currentLength) {
        if (maxLength < currentLength) {
          maxLength = currentLength;
        }
        //当前长度等于上次的长度减去（该字符上次出现位置减去新子串第一个元素的索引） 例如bbtablud遍历到第三个b实际减0
        currentLength -= intArray[idx] - (i - currentLength) - 1; //为了避免数组intArray初始化，采取+-1的方式还原真实索引
      } else {
        currentLength++;
      }
      intArray[idx] = i + 1;
    }
    return Math.max(maxLength, currentLength);
  }
}
