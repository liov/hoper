标题也许看起来有点绕，此时我的心情也有点遭，这是一个中等难度的题目，但是我却浪费了整整4天，并且结果还是失败的，不过好在在我的windows机子上测试是通过的。

我坚信这是rust的bug并不是我的，虽然我大量用了unsafe，但是情况确实windows上可以过去，playground跑步起来，leetcode测不过去。

如果是我的代码问题，那么三个平台的结论应该是一致，但显然不是这样。

早些时候我就发现，rust的程序在不同终端有不一致的表现，当然这仍然与unsafe有关。

但是，你是无法避免rust的unsafe的，至少是目前无法避免。尤其是rust的链表，没错这次我的问题刚刚好是在链表上。

LeetCode上的一道题：

给定一个链表，旋转链表，将链表每个节点向右移动 k 个位置，其中 k 是非负数。

# 示例 1:
```
输入: 1->2->3->4->5->NULL, k = 2
输出: 4->5->1->2->3->NULL
解释:
向右旋转 1 步: 5->1->2->3->4->NULL
向右旋转 2 步: 4->5->1->2->3->NULL
```
# 示例 2:
```
输入: 0->1->2->NULL, k = 4
输出: 2->0->1->NULL 解释: 向右旋转 1 步: 2->0->1->NULL 向右旋转 2 步: 1->2->0->NULL 向右旋转 3 步: 0->1->2->NULL 向右旋转 4 步: 2->0->1->NULL
起初我并不在意，但是你使用的是抽象0开销的语言，那么你的目标就是，不能增加内存消耗，也就是避免clone，我知道用clone并且不在乎时间复杂度，我是可以轻松写出来的，但这显然不是我要的。
```
起初为链表定义了一个求长度的方法，后来优化版发现是没用的，因为这样得遍历两边

```rust
impl ListNode {
    pub fn len(&self) -> usize {
        fn get_len(list: &ListNode, len: usize) -> usize {
            if let Some(ref next) = list.next {
                get_len(next, len + 1)
            } else {
                return len;
            }
        }
        get_len(self, 1)
    }
}
```
废话不多，上解题代码

```rust
pub fn rotate_right(head: Option<Box<ListNode>>, k: i32) -> Option<Box<ListNode>> {
    if k == 0 { return head; }
    let mut l = head.unwrap();

    fn get(list: &mut Box<ListNode>, len: i32, k: i32) -> (Option<*mut Box<ListNode>>, Option<*mut Box<ListNode>>, i32) {
        if let Some(ref mut next) = list.next {
            let (p2, p3, size) = get(next, len + 1, k);
            if len == size - k + 1 {
                return ( Some(list), p3, size);
            } else if len == size - k {
                let mut t =list.clone();
                t.next.take();
                mem::replace( list, t);
            }
            return (p2, p3, size);
        } else {
            return (None, Some(list), len);
        }

        panic!("错误")
    }

    let ( p2, p3, len) = get(&mut l, 1, k);

    if k > len {
        return rotate_right(Some(l), k % len);
    };

    unsafe {
        mem::replace(&mut (*p3.unwrap()).next, Some(l.clone()));
        Some((*p2.unwrap()).clone())
    }
}
```
测试代码，当然这之前得给出链表的定义，这是rust最简单的链表，正是因为简单才难处理

```rust
impl ListNode {
    #[inline]
    pub fn new(val: i32) -> Self {
        ListNode {
            next: None,
            val,
        }
    }

    /*因为 Rust 是后向兼容的（backwards compatible），所以不会移除 ref 和 ref mut，同时它们在一些不明确的场景还有用，
    比如希望可变地借用结构体的部分值而可变地借用另一部分的情况。你可能会在老的 Rust 代码中看到它们，所以请记住它们仍有价值。*/
    pub fn push(&mut self, val: i32) {
        match self.next {
            Some(ref mut next) =>
                next.push(val),
            None => self.next = Some(Box::new(ListNode::new(val)))
        }
    }
}
```
测试代码：
```rust
#[test]
fn rotate_right_test(){
    let mut l1= Box::new(ListNode::new(1));
    l1.push(2);
    l1.push(3);
    l1.push(3);
    l1.push(4);
    l1.push(5);
    let mut l2= Box::new(ListNode::new(3));
    l2.push(4);
    l2.push(5);
    l2.push(1);
    l2.push(2);
    l2.push(3);
    assert_eq!(rotate_right2(Some(l1),5), Some(l2))
}
```
windows上是测试通过的，不过写成的代码与我最初想象的完全不是一个东西，我设想的是组成环然后切开（当然实际是先切后组），避免任何一个clone

然而大量使用了unsafe，还用了clone，并且最终还无法通过LeetCode的测试

wtf？？？

以后在分析这段代码吧，我觉得我用了挺多心思的

但是rust的链表，链表，链表，链表，链表！！！！！