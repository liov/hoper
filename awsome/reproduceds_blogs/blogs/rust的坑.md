结合几日来用rust刷leetcode

发现rust的开发效率真的不敢恭维，当然这都与特定数据结构有关，可能因为自己不够深入吧

rust的所有权系统和某天数据结构天然的冲突，比如官方的链表是用unsafe实现的

rust无法无额外开销的将一个链表切开，并首尾相连组成新链表

rust的引用规则

在任意给定时间，要么 只能有一个可变引用，要么 只能有多个不可变引用。
引用必须总是有效。
因此你无法在遍历二叉树的时候，添加数据并重置迭代器

奇怪，强大的NLL

这段代码编译的过去，另一段不行
```rust
fn main() {
    let mut set: BTreeSet<usize> = [3, 1, 2].iter().cloned().collect();
    let mut set_iter = set.iter();
    let mut i = 0;
    while let Some(t) = set_iter.next() {
        println!("{:?}", set);
        i+=1;
        set.insert(i);
        set_iter = set.iter();
    }
}
```
问题来自LeetCode的接雨水
```rust
use std::collections::{HashSet,BTreeMap};

#[derive(PartialEq, Eq, PartialOrd, Ord, Debug, Hash, Clone, Copy)]
struct Point(i32, usize, usize);


pub fn trap_rain_water(height_map: Vec<Vec<i32>>) -> i32 {
    if height_map.len() < 3 || height_map[0].len() < 3 { return 0; };
    let m = height_map.len();
    let n = height_map[0].len();
    let mut result = 0;
    let mut side = BTreeMap::new();
    let mut drop = HashSet::new();
    for x in 0..m {
        for y in 0..n {
            if x == 0 || x == m - 1 || y == 0 || y == n - 1 {
                side.insert(Point(height_map[x][y], x, y),false);
                drop.insert((x,y));
            }
        }
    }

    let round: [[i32; 2]; 4] = [[0, -1], [0, 1], [-1, 0], [1, 0]];
    let mut x = 0;
    let mut y = 0;

    let mut point_iter = side.iter();
    while let Some((point,yet)) = point_iter.next() {
        if *yet { continue; };
        let point =*point;
        for j in round.iter() {
            x = (point.1 as i32 + j[0]) as usize;
            y = (point.2 as i32 + j[1]) as usize;
            if  x >= m - 1 || y >= n - 1 {
                continue;
            }
            if drop.get(&(x,y)) != None{
                continue;
            }
            drop.insert((x,y));
            if  height_map[x][y] <= point.0 {
                result = result + (point.0 -  height_map[x][y]);
                side.insert(Point(point.0, x, y),false);
                point_iter = side.iter();
            }else {
                side.insert(Point(height_map[x][y], x, y),false);
                point_iter = side.iter();
            }
        }
        side.insert(point,true);
        point_iter = side.iter();
    }
    result
}
```
用了个恶心的方法解决，每个insert后都重置迭代器