package leetcode;

import java.util.HashMap;
import java.util.HashSet;

public class CanFinish {
  public boolean canFinish(int numCourses, int[][] prerequisites) {
    var map = new HashMap<Integer, HashSet<Integer>>();
    for (int[] arr : prerequisites) {
      if (map.containsKey(arr[1])) {
        map.get(arr[1]).add(arr[0]);
      } else {
        var set = new HashSet<Integer>();
        set.add(arr[0]);
        map.put(arr[1], set);
      }
    }
    var set = new HashSet<Integer>();
    var mem = new HashSet<Integer>();
    for (int k : map.keySet()) {
      if (!dfs(k, set, mem, map)) return false;
    }
    return true;
  }

  public boolean dfs(int num, HashSet<Integer> set, HashSet<Integer> mem, HashMap<Integer, HashSet<Integer>> map) {
    if (set.contains(num)) return false;
    if (mem.contains(num)) return true;
    set.add(num);
    if (map.containsKey(num)) {
      for (int n : map.get(num)) {
        if (mem.contains(n)) return true;
        if (!dfs(n, set, mem, map)) return false;
      }
    }
    set.remove(num);
    return true;
  }
}
